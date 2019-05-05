package bruteforce

import (
	"github.com/xaionaro-go/algorithms/tsp/task"
)

type cache struct {
	cityAmount            uint
	cost                  []float64
	lastRoutesMinimalCost []float64
}

func newCache(t *task.Task) *cache {
	return &cache{
		cityAmount:            uint(len(t.Cities)),
		cost:                  make([]float64, len(t.Cities)*len(t.Cities)),
		lastRoutesMinimalCost: make([]float64, len(t.Cities)/2),
	}
}

func (c *cache) SetCost(startCityID, endCityID uint32, maxCost float64) {
	c.cost[int(startCityID)*int(c.cityAmount)+int(endCityID)] = maxCost
}

func (c *cache) GetCost(startCityID, endCityID uint32) float64 {
	return c.cost[int(startCityID)*int(c.cityAmount)+int(endCityID)]
}

func (c *cache) Prepare(w *worker, totalCostEstimation float64) {
	for _, city := range w.task.Cities {
		for _, route := range city.OutRoutes {
			_, cost := w.findSimplePath(
				route.StartCity,
				route.EndCity,
				1,
				nil,
				nil,
				0,
				route.Cost*0.999,
			)
			if cost > 0 && cost < route.Cost {
				//fmt.Println("found cheaper", route.StartCity.ID, route.EndCity.ID, cost, route.Cost, path)
				c.SetCost(route.StartCity.ID, route.EndCity.ID, cost)
			}
		}
	}

	startCity := w.task.StartCity
	for distanceMinusOne, _ := range c.lastRoutesMinimalCost {
		distance := distanceMinusOne + 1

		for _, city := range w.task.Cities {
			_, minimalCost := w.findCheapestPath(
				city,
				startCity,
				distance,
				nil,
				nil,
				0,
				nil,
				0,
				totalCostEstimation,
			)
			if c.lastRoutesMinimalCost[distanceMinusOne] == 0 || minimalCost < c.lastRoutesMinimalCost[distanceMinusOne] {
				//fmt.Println("minimal cost for distance", distance, "is", minimalCost)
				c.lastRoutesMinimalCost[distanceMinusOne] = minimalCost
			}
		}
	}
}
