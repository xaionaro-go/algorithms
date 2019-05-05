package bruteforce

import (
	"context"
	"math"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type worker struct {
	tick         int
	ctx          context.Context
	task         *task.Task
	intSlicePool *intSlicePool
	cache        *cache

	nonsetCityIDsTempBuffer [backScanDepth]uint32
}

func newWorker(ctx context.Context, t *task.Task) *worker {
	return &worker{
		ctx:          ctx,
		task:         t,
		intSlicePool: newIntSlicePool(),
		cache:        newCache(t),
	}
}

func (w *worker) isTimedOut() bool {
	if w.tick&0xffff == 0 {
		select {
		case <-w.ctx.Done(): // timeout
			return true
		default:
		}
	}
	return false
}

// Find any solution (but fast)
func (w *worker) findSimplePath(
	startCity *task.City,
	endCity *task.City,
	requireTotalCount int,
	cityCount []int,
	curPath *task.Path,
	curCost float64,
	costLimit float64,
) (task.Path, float64) {
	var city *task.City
	if curPath != nil {
		city = (*curPath)[len(*curPath)-1].EndCity
	} else {
		city = startCity
	}

	if curPath == nil {
		curPath = &task.Path{}
	}

	if len(*curPath) >= requireTotalCount && city == endCity {
		result := make(task.Path, len(*curPath))
		copy(result, *curPath)
		return result, curCost
	}

	w.tick++
	if w.isTimedOut() {
		return nil, math.Inf(-1)
	}

	if cityCount == nil {
		cityCount = make([]int, len(w.task.Cities))
	}

	for _, route := range city.OutRoutes {
		if cityCount[route.EndCity.ID] != 0 {
			continue // we already were in this city, skip it
		}
		if costLimit > 0 && curCost+route.Cost > costLimit {
			continue
		}

		cityCount[route.EndCity.ID]++
		*curPath = append(*curPath, route)
		path, cost := w.findSimplePath(
			startCity,
			endCity,
			requireTotalCount,
			cityCount,
			curPath,
			curCost+route.Cost,
			costLimit,
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

// Find the cheapest solution
func (w *worker) findCheapestPath(
	startCity *task.City,
	endCity *task.City,
	requireTotalCityCount int,
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
		city = startCity
	}

	//fmt.Println(cityCount, uselessCityCount, totalCityCount, t.StartCity.ID, costLimit, curPath)

	if totalCityCount >= requireTotalCityCount && (city == endCity || endCity == nil) {
		result := make(task.Path, len(*curPath))
		copy(result, *curPath)
		return result, curCost
	}

	w.tick++
	if w.isTimedOut() {
		return nil, math.Inf(-1)
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

	var minimalCostLeft float64
	if endCity == w.task.StartCity {
		distance := len(w.task.Cities) - totalCityCount
		if distance == backScanDepth {
			nonsetCityIdx := 0
			for cityID, count := range cityCount {
				if count == 0 {
					w.nonsetCityIDsTempBuffer[nonsetCityIdx] = uint32(cityID)
					nonsetCityIdx++
				}
			}

			costLeft := w.cache.GetLastRoutesCost(w.nonsetCityIDsTempBuffer[:])
			if costLeft > costLimit-curCost {
				return nil, -1
			}
		}
		distanceMinusOne := distance - 1
		if distanceMinusOne > 0 && len(w.cache.lastRoutesMinimalCost) > distanceMinusOne-1 {
			minimalCostLeft = w.cache.lastRoutesMinimalCost[distanceMinusOne-1]
		}
	}

	var cheapestPath task.Path
	var cheapestCost float64
	for _, route := range city.OutRoutes {
		if costLimit > 0 && curCost+route.Cost > costLimit {
			continue
		}
		cacheCost := w.cache.GetCost(route.StartCity.ID, route.EndCity.ID)
		if cacheCost > 0 && route.Cost > cacheCost {
			continue
		}
		if minimalCostLeft > 0 && costLimit-(curCost+route.Cost) < minimalCostLeft {
			continue
		}

		// To prevent loops we remember all cities that we revisit without visiting any new/unvisited cities
		// And if we return to the same (already visited) city without visiting any new cities, then it
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
			startCity,
			endCity,
			requireTotalCityCount,
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

func (w *worker) prepareCache(costEstimation float64) {
	w.cache.Prepare(w, costEstimation)
}
