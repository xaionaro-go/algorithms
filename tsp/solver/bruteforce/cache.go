package bruteforce

import (
	"github.com/xaionaro-go/algorithms/tsp/task"
)

type cache struct {
	cityAmount uint
	cost       []float64
}

func newCache(t *task.Task) *cache {
	return &cache{
		cityAmount: uint(len(t.Cities)),
		cost:       make([]float64, len(t.Cities)*len(t.Cities)),
	}
}

func (c *cache) SetCost(startCityID, endCityID uint32, maxCost float64) {
	c.cost[int(startCityID)*int(c.cityAmount)+int(endCityID)] = maxCost
}

func (c *cache) GetCost(startCityID, endCityID uint32) float64 {
	return c.cost[int(startCityID)*int(c.cityAmount)+int(endCityID)]
}
