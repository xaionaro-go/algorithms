package bruteforce

import (
	"context"
	"github.com/xaionaro-go/algorithms/tsp/task"
	"math"
	"math/rand"
)

type worker struct {
	tick         int
	ctx          context.Context
	task         *task.Task
	intSlicePool *intSlicePool
}

func newWorker(ctx context.Context, t *task.Task) *worker {
	return &worker{
		ctx:          ctx,
		task:         t,
		intSlicePool: newIntSlicePool(),
	}
}

func (w *worker) findSimplePath(
	cityCount []int,
	curPath *task.Path,
	curCost float64,
) (task.Path, float64) {
	var city *task.City
	if curPath != nil {
		city = (*curPath)[len(*curPath)-1].EndCity
	} else {
		city = w.task.StartCity
	}

	if curPath == nil {
		curPath = &task.Path{}
	}

	if len(*curPath) == len(w.task.Cities) && city == w.task.StartCity {
		result := make(task.Path, len(*curPath))
		copy(result, *curPath)
		return result, curCost
	}

	if rand.Intn(1000) == 0 {
		select {
		case <-w.ctx.Done(): // timeout
			return nil, math.Inf(-1)
		default:
		}
	}

	if cityCount == nil {
		cityCount = make([]int, len(w.task.Cities))
	}

	for _, route := range city.OutRoutes {
		if cityCount[route.EndCity.ID] != 0 {
			continue // we already were in this city, skip it
		}

		cityCount[route.EndCity.ID]++
		*curPath = append(*curPath, route)
		path, cost := w.findSimplePath(
			cityCount,
			curPath,
			curCost+route.Cost,
		)

		*curPath = (*curPath)[:len(*curPath)-1]
		cityCount[route.EndCity.ID]--

		if math.IsInf(cost, -1) { // timeout
			return path, cost
		}
		if cost < 0 { // a dead-end
			continue
		}
		if cost >= 0 {
			return path, cost // some working solution (not optimal, but correct), the end
		}
	}
	return nil, -1 // a dead-end
}

func (w *worker) findCheapestPath(
	cityCount []int,
	uselessCityCount *[]int,
	totalCityCount int,
	curPath *task.Path,
	curCost float64,
	costLimit float64,
) (task.Path, float64) {
	var city *task.City
	if curPath != nil {
		city = (*curPath)[len(*curPath)-1].EndCity
	} else {
		city = w.task.StartCity
	}

	//fmt.Println(cityCount, uselessCityCount, totalCityCount, t.StartCity.ID, costLimit, curPath)

	if totalCityCount == len(w.task.Cities) && city == w.task.StartCity {
		result := make(task.Path, len(*curPath))
		copy(result, *curPath)
		return result, curCost
	}

	w.tick++
	if w.tick&0xffff == 0 {
		select {
		case <-w.ctx.Done(): // timeout
			return nil, math.Inf(-1)
		default:
		}
	}

	if curPath == nil {
		curPath = &task.Path{}
	}

	if cityCount == nil {
		cityCount = make([]int, len(w.task.Cities))
	}

	if uselessCityCount == nil {
		uselessCityCount = &[][]int{make([]int, len(w.task.Cities))}[0]
	}

	var cheapestPath task.Path
	var cheapestCost float64
	for _, route := range city.OutRoutes {
		if costLimit > 0 && curCost+route.Cost > costLimit {
			continue
		}

		// To prevent loops we remember all cities that we revisit without visiting any new/unvisited cities
		// And if we return to the same (already visited) city without visiting any newre cities, then it
		// was an useless loop.
		if (*uselessCityCount)[route.EndCity.ID] != 0 {
			continue // there's no point to return to this city
		}

		newCostLimit := cheapestCost
		if newCostLimit <= 0 {
			newCostLimit = costLimit
		}

		var newUselessCityCount *[]int
		if cityCount[route.EndCity.ID] == 0 {
			totalCityCount++
			// An unvisited city, resetting "newUselessCityCount"
			newUselessCityCount = w.intSlicePool.Get(len(w.task.Cities), len(w.task.Cities))
		} else {
			// An already visited city
			(*uselessCityCount)[route.EndCity.ID]++
			newUselessCityCount = uselessCityCount
		}
		cityCount[route.EndCity.ID]++

		*curPath = append(*curPath, route)

		path, cost := w.findCheapestPath(
			cityCount,
			newUselessCityCount,
			totalCityCount,
			curPath,
			curCost+route.Cost,
			newCostLimit,
		)

		*curPath = (*curPath)[:len(*curPath)-1]

		cityCount[route.EndCity.ID]--
		if cityCount[route.EndCity.ID] == 0 {
			totalCityCount--
			w.intSlicePool.Put(newUselessCityCount)
		} else {
			(*uselessCityCount)[route.EndCity.ID]--
		}

		if math.IsInf(cost, -1) { // timeout
			return path, cost
		}
		if cost < 0 { // a dead-end
			continue
		}
		if cheapestPath == nil {
			cheapestCost = cost
			cheapestPath = path
		} else if cost < cheapestCost {
			cheapestCost = cost
			cheapestPath = path
		}
	}
	if cheapestPath == nil {
		return nil, -1 // a dead-end
	}

	return cheapestPath, cheapestCost
}
