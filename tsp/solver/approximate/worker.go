package approximate

import (
	"context"
	"fmt"
	"math"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

const (
	accuracyBacklash = 0.01
	bruteforceLength = 6
)

type worker struct {
	ctx          context.Context
	task         *task.Task
	tick         uint16
	cache        *cache
	intSlicePool *intSlicePool
}

func newWorker(ctx context.Context, t *task.Task) *worker {
	w := &worker{
		ctx:          ctx,
		task:         t,
		intSlicePool: newIntSlicePool(),
	}
	w.cache = newCache()
	w.cache.prepare(w, t)
	return w
}

func (w *worker) isTimedOut() bool {
	if w.tick == 0 {
		select {
		case <-w.ctx.Done(): // timeout
			return true
		default:
		}
	}
	return false
}

// Find the cheapest solution
func (w *worker) findCheapestPathFromCache(
	startCity *task.City,
	endCity *task.City,
	requireCityCount []int,
	cityCountLeft int,
	curPath *task.Path,
	curCost float64,
	costLimit float64,
) (task.Path, float64) {
	var city *task.City
	if len(*curPath) > 0 {
		city = (*curPath)[len(*curPath)-1].EndCity
	} else {
		city = startCity
	}

	//fmt.Println(cityCountLeft, w.task.StartCity.ID, curCost, costLimit, *curPath, endCity.ID, requireCityCount)

	if cityCountLeft == 1 {
		routePath, routeCost := w.cache.GetPath(city.ID, endCity.ID)
		if routePath == nil {
			return nil, -1 // a dead-end
		}
		result := make(task.Path, 0, len(*curPath)+len(routePath))
		result = append(result, *curPath...)
		result = append(result, routePath...)
		//fmt.Println("win", requireCityCount, curCost+routeCost, result)
		return result, curCost + routeCost
	}

	/*w.tick++
	if w.isTimedOut() {
		return nil, math.Inf(-1)
	}*/

	var cheapestPath task.Path
	var cheapestCost float64
	for cityID, requireCount := range requireCityCount {
		if requireCount <= 0 {
			continue
		}
		if city.ID == uint32(cityID) {
			continue
		}
		if uint32(cityID) == endCity.ID {
			continue
		}

		routePath, routeCost := w.cache.GetPath(city.ID, uint32(cityID))
		if routeCost <= 0 {
			continue
		}
		if costLimit > 0 && curCost+routeCost > costLimit*1.0001 {
			continue
		}

		*curPath = append(*curPath, routePath...)

		if len(cheapestPath) > 0 {
			matches := true
			for idx := range *curPath {
				if len(cheapestPath) < idx+1 {
					matches = false
					break
				}
				if (*curPath)[idx] != cheapestPath[idx] {
					matches = false
					break
				}
			}
			if matches {
				*curPath = (*curPath)[:len(*curPath)-len(routePath)]
				continue
			}
		}

		newCostLimit := cheapestCost
		if newCostLimit <= 0 {
			newCostLimit = costLimit
		}

		for _, route := range routePath {
			if route.EndCity.ID != endCity.ID {
				requireCityCount[route.EndCity.ID]--
			}
			if requireCityCount[route.EndCity.ID] == 0 {
				cityCountLeft--
			}
		}

		path, cost := w.findCheapestPathFromCache(
			startCity,
			endCity,
			requireCityCount,
			cityCountLeft,
			curPath,
			curCost+routeCost,
			newCostLimit,
		)

		*curPath = (*curPath)[:len(*curPath)-len(routePath)]

		for _, route := range routePath {
			if requireCityCount[route.EndCity.ID] == 0 {
				cityCountLeft++
			}
			if route.EndCity.ID != endCity.ID {
				requireCityCount[route.EndCity.ID]++
			}
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

// Find the cheapest solution
func (w *worker) findCheapestPath(
	startCity *task.City,
	endCity *task.City,
	requireCityCount []int,
	uselessCityCount *[]int,
	cityCountLeft int,
	curPath *task.Path,
	curCost float64,
	costLimit float64,
) (task.Path, float64) {
	var city *task.City
	if len(*curPath) > 0 {
		city = (*curPath)[len(*curPath)-1].EndCity
	} else {
		city = startCity
	}

	//fmt.Println(requireCityCount, uselessCityCount, cityCountLeft, w.task.StartCity.ID, costLimit, curPath)

	if cityCountLeft == 0 && (city == endCity || endCity == nil) {
		result := make(task.Path, len(*curPath))
		copy(result, *curPath)
		return result, curCost
	}

	/*w.tick++
	if w.isTimedOut() {
		return nil, math.Inf(-1)
	}*/

	var cheapestPath task.Path
	var cheapestCost float64
	for _, route := range city.OutRoutes {
		if costLimit > 0 && curCost+route.Cost > costLimit {
			continue
		}
		_, cacheCost := w.cache.GetPath(route.StartCity.ID, route.EndCity.ID)
		if cacheCost > 0 && route.Cost > cacheCost {
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

		newCity := requireCityCount[route.EndCity.ID] > 0
		var newUselessCityCount *[]int
		if newCity {
			// An unvisited city, resetting "newUselessCityCount"
			newUselessCityCount = w.intSlicePool.Get(len(w.task.Cities), len(w.task.Cities))
			requireCityCount[route.EndCity.ID]--
			if requireCityCount[route.EndCity.ID] == 0 {
				cityCountLeft--
			}
		} else {
			// An already visited city
			(*uselessCityCount)[route.EndCity.ID]++
			newUselessCityCount = uselessCityCount
		}

		*curPath = append(*curPath, route)

		path, cost := w.findCheapestPath(
			startCity,
			endCity,
			requireCityCount,
			newUselessCityCount,
			cityCountLeft,
			curPath,
			curCost+route.Cost,
			newCostLimit,
		)

		*curPath = (*curPath)[:len(*curPath)-1]

		if newCity {
			if requireCityCount[route.EndCity.ID] == 0 {
				cityCountLeft++
			}
			requireCityCount[route.EndCity.ID]++
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

// Find any solution (but fast)
func (w *worker) findSimplePath(
	startCity *task.City,
	endCity *task.City,
	requiredCityCount []int,
	cityCountLeft int,
	curPath *task.Path,
	curCost float64,
	costLimit float64,
) (task.Path, float64) {
	var city *task.City
	if len(*curPath) > 0 {
		city = (*curPath)[len(*curPath)-1].EndCity
	} else {
		city = startCity
	}

	//fmt.Println(startCity.ID, endCity.ID, requiredCityCount, cityCountLeft, *curPath, curCost, costLimit)

	if cityCountLeft <= 0 && city == endCity {
		result := make(task.Path, len(*curPath))
		copy(result, *curPath)
		return result, curCost
	}

	/*
		c.tick++
		if c.isTimedOut() {
			return nil, math.Inf(-1)
		}
	*/

	if requiredCityCount == nil {
		panic(`wrong arguments`)
	}

	for _, route := range city.OutRoutes {
		if requiredCityCount[route.EndCity.ID] == 0 {
			continue // we already were in this city, skip it
		}
		if costLimit > 0 && curCost+route.Cost > costLimit {
			continue
		}

		requiredCityCount[route.EndCity.ID]--
		*curPath = append(*curPath, route)
		path, cost := w.findSimplePath(
			startCity,
			endCity,
			requiredCityCount,
			cityCountLeft-1,
			curPath,
			curCost+route.Cost,
			costLimit,
		)

		*curPath = (*curPath)[:len(*curPath)-1]
		requiredCityCount[route.EndCity.ID]++

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

func (w *worker) findPath(
	startCity *task.City,
	endCity *task.City,
	requiredCityCountOrig []int,
) (task.Path, float64) {
	requiredCityCount := make([]int, len(requiredCityCountOrig))
	copy(requiredCityCount, requiredCityCountOrig)

	citiesLeft := 0
	for _, countLeft := range requiredCityCount {
		if countLeft > 0 {
			citiesLeft++
		}
	}

	curCity := startCity
	var path task.Path
	for citiesLeft > 0 {
		bestScore := -math.MaxFloat64
		var bestScorePath task.Path
		for cityID, countLeft := range requiredCityCount {
			if countLeft <= 0 {
				continue
			}
			if citiesLeft > 1 && uint32(cityID) == endCity.ID {
				continue
			}
			if uint32(cityID) == curCity.ID {
				continue
			}
			path, costToCity := w.cache.GetPath(curCity.ID, uint32(cityID))
			if costToCity <= 0 {
				continue
			}
			_, costFromCityToEnd := w.cache.GetPath(uint32(cityID), endCity.ID)
			if costFromCityToEnd <= 0 {
				costFromCityToEnd = w.cache.GetMaxCost()
			}
			score := (costFromCityToEnd + float64(citiesLeft-1)/float64(len(w.task.Cities))*w.cache.GetAverageCost()) / costToCity
			score = costFromCityToEnd / costToCity
			//score = 1 / costToCity
			if score > bestScore {
				bestScore = score
				bestScorePath = path
			}
		}

		if bestScore < 0 || len(bestScorePath) == 0 {
			panic(fmt.Sprintln(`Unexpected situation`, path, curCity.ID, requiredCityCount, bestScore, bestScorePath))
			uselessCityCount := make([]int, len(w.task.Cities))
			path, cost := w.findCheapestPath(
				startCity,
				endCity,
				requiredCityCount,
				&uselessCityCount,
				citiesLeft,
				&path,
				path.Cost(),
				0,
			)
			fmt.Println("bruteforce", path, cost, citiesLeft, requiredCityCount)
			return path, cost
		}

		path = append(path, bestScorePath...)
		for _, route := range bestScorePath {
			if citiesLeft == 1 || route.EndCity.ID != endCity.ID {
				requiredCityCount[route.EndCity.ID]--
				if requiredCityCount[route.EndCity.ID] == 0 {
					citiesLeft--
				}
			}
		}

		curCity = bestScorePath[len(bestScorePath)-1].EndCity
	}

	if curCity != endCity {
		lastMile, _ := w.cache.GetPath(curCity.ID, endCity.ID)
		path = append(path, lastMile...)
	}

	return path, path.Cost()
}

func (w *worker) optimizePath(path task.Path) task.Path {
	//fmt.Println("optimizePath", path)
	cityLeftCount := make([]int, len(w.task.Cities))
	cityLeftCountForBruteforce := make([]int, len(w.task.Cities))
	pathTmp := make(task.Path, len(w.task.Cities)*3)

	oldCost := math.MaxFloat64
	startIdx := 0
optimizePathLoop:
	for path.Cost() < oldCost || startIdx != 0 {
		if path[0].StartCity.ID != w.task.StartCity.ID {
			panic(fmt.Sprintln(w.task.StartCity.ID, path, oldCost, startIdx))
		}
		oldCost = path.Cost()
		for idx := range cityLeftCount {
			cityLeftCount[idx] = 1
			cityLeftCountForBruteforce[idx] = 1
		}
		curSegmentCost := float64(0)
		for idx, route := range path {
			if idx > bruteforceLength {
				cityLeftCountForBruteforce[path[idx-bruteforceLength-1].EndCity.ID]++
			}
			endCity := route.EndCity
			if cityLeftCount[route.EndCity.ID] == 1 {
				cityLeftCount[route.EndCity.ID]--
			}
			cityLeftCountForBruteforce[route.EndCity.ID]--
			if idx < startIdx {
				continue
			}
			curSegmentCost += route.Cost

			cityCount := make([]int, len(w.task.Cities))
			for cityID := range cityCount {
				cityCount[cityID] = 1 - cityLeftCount[cityID]
			}

			otherPath, otherSegmentCost := w.findPath(w.task.StartCity, endCity, cityCount)
			if otherSegmentCost > 0 && otherSegmentCost < curSegmentCost*0.9999 {
				//fmt.Println("optimized: ", path[:idx+1], otherPath, curSegmentCost, path, cityCount)
				path = append(otherPath, path[idx+1:]...)
				if path.Cost() > oldCost {
					panic(fmt.Sprintln(`Shouldn't happened'`, path, oldCost))
				}
				startIdx = len(otherPath)
				//fmt.Println("new path", path, oldCost, startIdx)
				continue optimizePathLoop
			}

			idxS := idx - bruteforceLength
			if idxS < 0 {
				idxS = 0
			}
			if idxS == idx {
				continue
			}

			subPath := path[idxS : idx+1]
			subPathCost := subPath.Cost()

			cityCountCount := 0
			for cityID := range cityCount {
				if cityLeftCountForBruteforce[cityID] > 0 {
					cityCount[cityID] = 0
				} else {
					cityCount[cityID] = 1
					cityCountCount++
				}
			}

			//fmt.Println("try bf", idxS, idx, subPath, cityCount, cityLeftCountForBruteforce)
			pathTmp = pathTmp[:0]
			otherPath, otherSegmentCost = w.findCheapestPathFromCache(
				subPath[0].StartCity,
				subPath[len(subPath)-1].EndCity,
				cityCount,
				cityCountCount,
				&pathTmp,
				0,
				subPathCost,
			)
			if otherSegmentCost <= 0 {
				continue
			}
			if otherSegmentCost >= subPathCost*0.9999 {
				continue
			}

			//fmt.Println("bf optimized: ", idxS, idx, subPath, otherPath, path, cityCount, path[:idxS], otherPath, path[idx+1:], path[idxS:idx+1])
			newPath := make(task.Path, 0, idxS+len(otherPath)+len(path)-idx)
			newPath = append(newPath, path[:idxS]...)
			newPath = append(newPath, otherPath...)
			newPath = append(newPath, path[idx+1:]...)
			path = newPath
			if path.Cost() > oldCost {
				panic(fmt.Sprintln(`Shouldn't happened'`, path, oldCost))
			}
			startIdx = len(otherPath)
			//fmt.Println("new path", path, oldCost, startIdx)
			continue optimizePathLoop
		}

		startIdx = 0
	}

	return path
}
