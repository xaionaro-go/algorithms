package approximate

import (
	"fmt"
	"github.com/xaionaro-go/algorithms/tsp/task"
	"math"
	"sort"
)

type cache struct {
	cityAmount           uint
	cost                 []float64
	path                 []task.Path
	maxCost              float64
	avgCost              float64
	routeEffectiveness   []int
	minimalLastRouteCost float64
}

func newCache() *cache {
	return &cache{}
}

func (c *cache) setAverageCost(newAvgCost float64) {
	c.avgCost = newAvgCost
}

func (c *cache) GetAverageCost() float64 {
	return c.avgCost
}

func (c *cache) setMaxCost(newMaxCost float64) {
	c.maxCost = newMaxCost
}

func (c *cache) GetMaxCost() float64 {
	return c.maxCost
}

func (c *cache) setPath(startCityID, endCityID uint32, path task.Path, cost float64) {
	key := int(startCityID)*int(c.cityAmount) + int(endCityID)
	c.path[key] = path
	c.cost[key] = cost
	//fmt.Println(startCityID, endCityID, path)
	if cost > c.GetMaxCost() {
		c.setMaxCost(cost)
	}
}

func (c *cache) incrementRouteEffectiveness(startCityID, endCityID uint32) {
	key := int(startCityID)*int(c.cityAmount) + int(endCityID)
	c.routeEffectiveness[key]++
}

func (c *cache) GetRouteEffectiveness(startCityID, endCityID uint32) int {
	key := int(startCityID)*int(c.cityAmount) + int(endCityID)
	return c.routeEffectiveness[key]
}

func (c *cache) GetPath(startCityID, endCityID uint32) (task.Path, float64) {
	key := int(startCityID)*int(c.cityAmount) + int(endCityID)
	return c.path[key], c.cost[key]
}

type candidate struct {
	destination  *task.City
	cost         float64
	path         task.Path
	isEstimation bool
}
type candidates []*candidate

func (s candidates) Sort() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].cost < s[j].cost
	})
}

func (c *cache) sortTaskData(minimalInRouteCost []float64, t *task.Task) {
	for _, city := range t.Cities {
		sort.Slice(city.OutRoutes, func(i, j int) bool {
			aRoute := city.OutRoutes[i]
			bRoute := city.OutRoutes[j]
			aCity := t.Cities[aRoute.EndCity.ID]
			bCity := t.Cities[bRoute.EndCity.ID]
			aRouteEffectiveness := float64(c.GetRouteEffectiveness(aRoute.StartCity.ID, aRoute.EndCity.ID)) / 100
			bRouteEffectiveness := float64(c.GetRouteEffectiveness(bRoute.StartCity.ID, bRoute.EndCity.ID)) / 100
			aScore := (aRouteEffectiveness + 1 + minimalInRouteCost[aCity.ID]/aRoute.Cost) / aRoute.Cost / aRoute.Cost
			bScore := (bRouteEffectiveness + 1 + minimalInRouteCost[bCity.ID]/bRoute.Cost) / bRoute.Cost / bRoute.Cost
			return aScore > bScore
		})
	}
}

func (c *cache) sortRoutesByDestinationAndCost(t *task.Task) {
	for _, city := range t.Cities {
		sort.Slice(city.OutRoutes, func(i, j int) bool {
			aRoute := city.OutRoutes[i]
			bRoute := city.OutRoutes[j]
			aCity := t.Cities[aRoute.EndCity.ID]
			bCity := t.Cities[bRoute.EndCity.ID]
			if aCity.ID != bCity.ID {
				return aCity.ID < bCity.ID
			}
			return aRoute.Cost < bRoute.Cost
		})
	}
}

func (c *cache) initialSortData(minimalInRouteCost []float64, t *task.Task) {
	for _, city := range t.Cities {
		sort.Slice(city.OutRoutes, func(i, j int) bool {
			aRoute := city.OutRoutes[i]
			bRoute := city.OutRoutes[j]
			aCity := t.Cities[aRoute.EndCity.ID]
			bCity := t.Cities[bRoute.EndCity.ID]
			aScore := (1 + minimalInRouteCost[aCity.ID]/aRoute.Cost) / aRoute.Cost / aRoute.Cost / (1 + math.Log(float64(len(aCity.InRoutes))))
			bScore := (1 + minimalInRouteCost[bCity.ID]/bRoute.Cost) / bRoute.Cost / bRoute.Cost / (1 + math.Log(float64(len(bCity.InRoutes))))
			return aScore > bScore
		})
	}
}

func (c *cache) prepare(w *worker, t *task.Task) {
	c.sortRoutesByDestinationAndCost(t)

	for _, city := range t.Cities {
		var newRoutes task.Routes
		cmpRoute := city.OutRoutes[0]
		newRoutes = append(newRoutes, cmpRoute)
		for _, route := range city.OutRoutes[1:] {
			if route.EndCity.ID == cmpRoute.EndCity.ID {
				continue
			}
			newRoutes = append(newRoutes, route)
			cmpRoute = route
		}
		city.OutRoutes = newRoutes
	}

	minimalInRouteCost := make([]float64, len(t.Cities))
	for _, city := range t.Cities {
		min := city.InRoutes[0].Cost
		for _, route := range city.InRoutes {
			if min > route.Cost {
				min = route.Cost
			}
		}
		minimalInRouteCost[city.ID] = min
	}

	c.initialSortData(minimalInRouteCost, t)

	c.cityAmount = uint(len(t.Cities))
	c.cost = make([]float64, c.cityAmount*c.cityAmount)
	c.path = make([]task.Path, c.cityAmount*c.cityAmount)
	c.routeEffectiveness = make([]int, c.cityAmount*c.cityAmount)

	costSumLog := float64(0)
	costSum := float64(0)
	count := 0

	cityCountAll := make([]int, len(t.Cities))
	for idx := range cityCountAll {
		cityCountAll[idx] = 1
	}
	uselessCityCount := make([]int, len(t.Cities))
	requireCityCount := make([]int, len(t.Cities))
	pathTmp := make(task.Path, len(w.task.Cities)*3)
	candidates := make(candidates, len(t.Cities))
	for idx := range candidates {
		candidates[idx] = &candidate{}
	}
	prevSortCount := 0
	for _, cityA := range t.Cities {
		for idx, cityB := range t.Cities {
			candidate := candidates[idx]
			candidate.destination = nil
			candidate.cost = math.MaxFloat64
			candidate.path = nil
			if cityA.ID == cityB.ID {
				continue
			}
			if path, cost := c.GetPath(cityA.ID, cityB.ID); path != nil {
				costSum += cost
				costSumLog += math.Log(cost)
				count++
				continue
			}

			candidate.isEstimation = true
			requireCityCount[cityB.ID] = 1
			pathTmp = pathTmp[:0]
			estimationPath, estimationCost := w.findSimplePath(
				cityA,
				cityB,
				cityCountAll,
				1,
				&pathTmp,
				0,
				0,
				5,
			)

			candidate.destination = cityB
			candidate.cost = estimationCost
			candidate.path = estimationPath

			if estimationPath == nil {
				panic(fmt.Sprintln(`Should not happened`, cityA.ID, cityB.ID, requireCityCount))
				/*candidate.isEstimation = false
				pathTmp = pathTmp[:0]
				path, cost := w.findCheapestPath(
					cityA,
					cityB,
					requireCityCount,
					&uselessCityCount,
					1,
					&pathTmp,
					0,
					estimationCost,
				)
				candidate.cost = cost
				candidate.path = path

				if path == nil {
					panic(`Shouldn't happened`)
				}*/
			}

			requireCityCount[cityB.ID] = 0
		}

		//candidates.Sort()

		for idx, candidate := range candidates {
			cityB := candidate.destination
			cost := candidate.cost
			path := candidate.path

			if cityB == nil {
				continue
			}

			if candidate.isEstimation {
				costThresholdTemp := float64(0)
				for _, route := range cityA.OutRoutes {
					costThresholdTemp += math.Log(route.Cost)
				}
				costThreshold := math.Exp(costThresholdTemp / float64(len(cityA.OutRoutes)))
				if costThreshold > math.Exp(costSumLog/float64(count)) {
					costThreshold = math.Exp(costSumLog / float64(count))
				}

				if cityB.ID == t.StartCity.ID || (count < 100 && (idx < 10 || cost < costThreshold)) {
					requireCityCount[cityB.ID] = 1
					pathTmp = pathTmp[:0]
					path, cost = w.findCheapestPath(
						cityA,
						cityB,
						requireCityCount,
						&uselessCityCount,
						1,
						&pathTmp,
						0,
						cost,
					)
					requireCityCount[cityB.ID] = 0
					candidate.isEstimation = false
				}
			}
			costSum += cost
			costSumLog += math.Log(cost)
			count++

			for idxS, routeS := range path {
				subPathCost := float64(0)
				c.incrementRouteEffectiveness(routeS.StartCity.ID, routeS.EndCity.ID)
				for idxE, routeE := range path[idxS:] {
					subPath := path[idxS : idxS+idxE+1]
					subPathCost += routeE.Cost
					oldPath, oldCost := c.GetPath(routeS.StartCity.ID, routeE.EndCity.ID)
					if oldPath != nil && oldCost <= subPathCost {
						continue
					}
					c.setPath(routeS.StartCity.ID, routeE.EndCity.ID, subPath, subPathCost)
				}
			}
		}

		if count > 10 && prevSortCount*4/3 < count {
			c.sortTaskData(minimalInRouteCost, t)
			prevSortCount = count
		}
	}
	c.sortTaskData(minimalInRouteCost, t)

	c.setAverageCost(costSum / float64(c.cityAmount*c.cityAmount))
}
