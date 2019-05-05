package bruteforce

import (
	"github.com/xaionaro-go/algorithms/tsp/task"
)

const (
	backScanDepth = 0
)

type cache struct {
	cityAmount               uint
	cost                     []float64
	lastRoutesMinimalCost    []float64
	couldBeUsedForWalkAround []bool
	lastRoutesCost           []float64

	cityIDsTempBuffer [backScanDepth]uint32
}

func pow(x, p int) int {
	r := 1
	for i := 0; i < p; i++ {
		r *= x
	}
	return r
}

func newCache(t *task.Task) *cache {
	return &cache{
		cityAmount:               uint(len(t.Cities)),
		cost:                     make([]float64, len(t.Cities)*len(t.Cities)),
		lastRoutesMinimalCost:    make([]float64, len(t.Cities)/2),
		lastRoutesCost:           make([]float64, pow(len(t.Cities)+1, backScanDepth)),
		couldBeUsedForWalkAround: make([]bool, len(t.Cities)*len(t.Cities)),
	}
}

type uint32Slice []uint32

func (s uint32Slice) Len() int           { return len(s) }
func (s uint32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s uint32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// copied from https://github.com/demdxx/sort-algorithms/blob/master/algorithms.go
func bubbleSort(data uint32Slice) {
	n := data.Len() - 1
	b := false
	for i := 0; i < n; i++ {
		for j := 0; j < n-i; j++ {
			if data.Less(j+1, j) {
				data.Swap(j+1, j)
				b = true
			}
		}
		if !b {
			break
		}
		b = false
	}
}

func (c *cache) getPathIDByCityIDs(cityIDs []uint32) uint64 {
	pathID := uint64(0)
	bubbleSort(uint32Slice(cityIDs))
	for _, cityID := range c.cityIDsTempBuffer {
		pathID *= uint64(c.cityAmount)
		pathID += uint64(cityID)
	}
	return pathID
}

func (c *cache) SetCost(startCityID, endCityID uint32, maxCost float64) {
	c.cost[int(startCityID)*int(c.cityAmount)+int(endCityID)] = maxCost
}

func (c *cache) GetCost(startCityID, endCityID uint32) float64 {
	return c.cost[int(startCityID)*int(c.cityAmount)+int(endCityID)]
}

func (c *cache) GetLastRoutesCost(nonsetCities []uint32) float64 {
	return c.lastRoutesCost[c.getPathIDByCityIDs(nonsetCities)]
}

func (c *cache) SetCouldBeUsedForWalkAround(startCityID, endCityID uint32) {
	c.couldBeUsedForWalkAround[startCityID*uint32(c.cityAmount)+endCityID] = true
}

func (c *cache) GetCouldBeUsedForWalkAround(startCityID, endCityID uint32) bool {
	return c.couldBeUsedForWalkAround[startCityID*uint32(c.cityAmount)+endCityID]
}

func (c *cache) cityCombinationsOfPath(path task.Path, size uint, fn func(cityIDs []uint32)) {
	cityIDs := make([]uint32, size)

	if uint(len(path)) < size {
		return
	}

	if uint(len(path)) == size {
		for idx, route := range path {
			cityIDs[idx] = route.StartCity.ID
		}
		fn(cityIDs)
		return
	}

	b := make([]bool, len(path))
	for i := uint(0); i < size; i++ {
		b[i] = true
	}

	for {
		cityIDsIdx := 0
		for idx, route := range path {
			if !b[idx] {
				continue
			}
			cityIDs[cityIDsIdx] = route.StartCity.ID
			cityIDsIdx++
		}
		fn(cityIDs[:])

		finished := true
		for i := uint(0); i < size; i++ {
			if !b[uint(len(path))-i-1] {
				finished = false
				break
			}
		}
		if finished {
			break
		}

		idx := len(path) - 1
		for ; idx >= 0 && b[idx]; idx-- {
		}

		for ; idx >= 0 && !b[idx]; idx-- {
		}

		b[idx+1] = true
		b[idx] = false

		rIdx := idx + 2
		for i := len(path) - 1; i > rIdx && b[i]; i-- {
			for ; b[rIdx] && rIdx < i; rIdx++ {
			}
			if rIdx >= i {
				break
			}
			b[rIdx] = true
			b[i] = false
			rIdx++
		}
	}
}

func (c *cache) Prepare(w *worker, totalCostEstimation float64) {
	for _, city := range w.task.Cities {
		for _, route := range city.OutRoutes {
			_, cost := w.findCheapestPath(
				route.StartCity,
				route.EndCity,
				1,
				nil,
				nil,
				0,
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

	for _, cityA := range w.task.Cities {
		for _, cityB := range w.task.Cities {
			_, estimationCost := w.findSimplePath(
				cityA,
				cityB,
				1,
				nil,
				nil,
				0,
				totalCostEstimation,
			)
			path, _ := w.findCheapestPath(
				cityA,
				cityB,
				1,
				nil,
				nil,
				0,
				nil,
				0,
				estimationCost,
			)
			if len(path) > 1 {
				for _, route := range path {
					c.SetCouldBeUsedForWalkAround(route.StartCity.ID, route.EndCity.ID)
				}
			}
		}
	}

	startCity := w.task.StartCity
	for distanceMinusOne, _ := range c.lastRoutesMinimalCost {
		distance := distanceMinusOne + 1

		for _, city := range w.task.Cities {
			path, minimalCost := w.findCheapestPath(
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
				c.lastRoutesMinimalCost[distanceMinusOne] = minimalCost
			}
			if distance == backScanDepth {
				c.cityCombinationsOfPath(path, backScanDepth, func(cityIDs []uint32) {
					pathID := c.getPathIDByCityIDs(cityIDs)
					if c.lastRoutesCost[pathID] == 0 || minimalCost < c.lastRoutesCost[pathID] {
						c.lastRoutesCost[pathID] = minimalCost
					}
				})
			}
		}
	}
}

func init() {

}
