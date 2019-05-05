package bruteforce

import (
	"math"
	"math/rand"
	"time"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type subWorker struct {
	*worker
	id               uint32
	tick             uint64
	cachedLeaderCost float64
}

func (w *subWorker) getLeaderCost() float64 {
	if w.tick&0x1ff == 0 || w.cachedLeaderCost == 0 {
		w.cachedLeaderCost = w.worker.leaderCost.Get()
	}
	return w.cachedLeaderCost
}

// Find any solution (but fast)
func (w *subWorker) findSimplePath(args *jobArguments, result *jobResult) {
	var city *task.City
	if len(args.curPath) > 0 {
		city = args.curPath[len(args.curPath)-1].EndCity
	} else {
		city = w.task.StartCity
	}

	if len(args.curPath) == len(w.task.Cities) && city == w.task.StartCity {
		w.considerSolution(w.id, args, result)
		return
	}

	if !w.isRunning {
		if result.cost == 0 {
			result.PointTo(jobResultCancel)
		}
		return
	}

	if args.cityCount == nil {
		args.cityCount = make([]int, len(w.task.Cities))
	}

	for _, route := range city.OutRoutes {
		if args.cityCount[route.EndCity.ID] != 0 {
			continue // we already were in this city, skip it
		}

		args.cityCount[route.EndCity.ID]++
		args.curPath = append(args.curPath, route)

		oldCurCost := args.curCost
		args.curCost += route.Cost

		w.findSimplePath(args, result)
		cost := result.cost

		args.curCost = oldCurCost

		args.curPath = args.curPath[:len(args.curPath)-1]
		args.cityCount[route.EndCity.ID]--

		if math.IsInf(cost, -1) { // timeout
			result.PointTo(jobResultCancel)
			return
		}
		if cost < 0 { // a dead-end
			continue
		}
		if cost >= 0 {
			return // some working solution (not optimal, but correct), the end
		}
	}
	result.PointTo(jobResultDeadEnd)
	return
}

func (w *subWorker) enqueue(args *jobArguments, result *jobResult) *jobResult {
	w.tick++
	if w.queue.Length()+w.parallelFactor >= w.queue.Size() || w.parallelFactor == 1 || w.tick&0xfffffff00 != 0 {
		w.callKernelFunc(args, result)
		return nil
	}

	job := w.jobPool.Get()
	job.args = args.Clone(w.jobArgumentsPool)
	jobResult := w.jobResultPool.Get()
	job.result = jobResult
	if !w.queue.Enqueue(job) {
		w.callKernelFunc(args, result)
		job.args.Release(w.jobArgumentsPool)
		job.result.Release(w.jobResultPool)
		job.args = nil
		job.result = nil
		job.Release(w.jobPool)
		return nil
	}
	return jobResult

}

// Find the cheapest solution
func (w *subWorker) findCheapestPath(args *jobArguments, result *jobResult) {
	var city *task.City
	if len(args.curPath) > 0 {
		city = args.curPath[len(args.curPath)-1].EndCity
	} else {
		city = w.task.StartCity
	}

	//fmt.Println("job", subWorkerID, args, w.getGlobalCostLimit(), args.curPath)

	if args.totalCityCount == len(w.task.Cities) && city == w.task.StartCity {
		w.considerSolution(w.id, args, result)
		return
	}

	if !w.isRunning {
		if result.cost == 0 {
			result.PointTo(jobResultCancel)
		}
		return
	}

	var processLater []*jobResult
	var cheapestResult *jobResult
	for _, route := range city.OutRoutes {
		if args.costLimit > 0 && args.curCost+route.Cost > args.costLimit {
			continue
		}
		if args.curCost+route.Cost > w.getLeaderCost() {
			continue
		}

		// To prevent loops we remember all cities that we revisit without visiting any new/unvisited cities
		// And if we return to the same (already visited) city without visiting any newre cities, then it
		// was an useless loop.
		if args.uselessCityCount[route.EndCity.ID] != 0 {
			continue // there's no point to return to this city
		}

		newCostLimit := args.costLimit
		if cheapestResult != nil {
			if cheapestResult.cost < newCostLimit {
				newCostLimit = cheapestResult.cost
			}
		}

		oldUselessCityCount := args.uselessCityCount

		var newUselessCityCount *[]int
		if args.cityCount[route.EndCity.ID] == 0 {
			args.totalCityCount++
			// An unvisited city, resetting "newUselessCityCount"
			newUselessCityCount = w.intSlicePool.Get(len(w.task.Cities), len(w.task.Cities))
			args.uselessCityCount = *newUselessCityCount
		} else {
			// An already visited city
			args.uselessCityCount[route.EndCity.ID]++
		}
		args.cityCount[route.EndCity.ID]++

		args.curPath = append(args.curPath, route)

		oldCurCost := args.curCost
		oldCostLimit := args.costLimit

		args.curCost += route.Cost
		args.costLimit = newCostLimit

		pendingResult := w.enqueue(args, result)

		args.uselessCityCount = oldUselessCityCount
		args.curCost = oldCurCost
		args.costLimit = oldCostLimit

		args.curPath = args.curPath[:len(args.curPath)-1]

		args.cityCount[route.EndCity.ID]--
		if args.cityCount[route.EndCity.ID] == 0 {
			args.totalCityCount--
			w.intSlicePool.Put(newUselessCityCount)
		} else {
			args.uselessCityCount[route.EndCity.ID]--
		}

		if pendingResult != nil {
			//fmt.Printf("%v waiting: %p\n", subWorkerID, pendingResult)
			processLater = append(processLater, pendingResult)
			continue
		}

		if math.IsInf(result.cost, -1) { // timeout
			return
		}
		if result.cost < 0 { // a dead-end
			continue
		}
		if cheapestResult == nil {
			cheapestResult = w.jobResultPool.Get()
			cheapestResult.ready = true
			cheapestResult.CopyFrom(result)
		} else if result.cost < cheapestResult.cost {
			cheapestResult.CopyFrom(result)
		}
	}
	if len(processLater) > 0 {
		for _, pendingResult := range processLater {
			for !pendingResult.ready {
				if !w.isRunning {
					if result.cost == 0 {
						result.PointTo(jobResultCancel)
					}
					return
				}
				w.completeTheRestWork()
			}
			if math.IsInf(pendingResult.cost, -1) { // timeout
				result.PointTo(jobResultCancel)
				pendingResult.Release(w.jobResultPool)
				return
			}
			if pendingResult.cost < 0 { // a dead-end
				pendingResult.Release(w.jobResultPool)
				continue
			}
			if cheapestResult == nil {
				cheapestResult = pendingResult
			} else if result.cost < cheapestResult.cost {
				cheapestResult.Release(w.jobResultPool)
				cheapestResult = pendingResult
			} else {
				pendingResult.Release(w.jobResultPool)
			}
		}
	}

	if cheapestResult == nil {
		result.PointTo(jobResultDeadEnd) // a dead-end
		return
	}

	result.CopyFrom(cheapestResult)
	cheapestResult.Release(w.jobResultPool)
	return
}

func (w *subWorker) callKernelFunc(args *jobArguments, result *jobResult) {
	switch w.mode {
	case workerModeSimplePath:
		w.findSimplePath(args, result)
	case workerModeFull:
		w.findCheapestPath(args, result)
	}
}

func (w *subWorker) doJob(job *job) {
	result := job.result
	job.result = nil

	w.callKernelFunc(job.args, result)
	if result.cost > 0 && (job.args.curPath == nil || len(job.args.curPath) == 0) {
		w.stop()
	}
	job.args.Release(w.jobArgumentsPool)
	job.args = nil
	job.Release(w.jobPool)
	//job.Release(w.jobPool, w.jobArgumentsPool)
	//fmt.Printf("%p %v %v %v\n", job.result, job.result.path, job.result.cost, job.result.ready)
	result.ready = true
}

func (w *subWorker) completeTheRestWork() {
	left := w.queue.Length()
	if left == 0 {
		return
	}
	w.subWorkersCount.Add(1)

	for i := uint32(0); i < left; i++ {
		job := w.queue.Dequeue()
		if job == nil {
			break
		}
		w.doJob(job)
	}

	w.subWorkersCount.Add(-1)
}

func (w *subWorker) execute() {
	for w.isRunning {
		w.completeTheRestWork()

		/*w.tick++
		if w.tick & 0xfff == 0 {
			fmt.Println(w.isRunning, w.subWorkersCount.Get(), w.queue.Length())
		}*/
		if w.subWorkersCount.Get() == 0 && w.queue.Length() == 0 {
			return
		}
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Nanosecond)
	}
}
