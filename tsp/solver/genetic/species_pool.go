package genetic

import (
	"sync"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type SpeciesPool struct {
	*sync.Pool
}

func newSpeciesPool(t *task.Task) *SpeciesPool {
	speciesPool := &SpeciesPool{&sync.Pool{}}
	speciesPool.New = func() interface{} {
		species := &Species{
			task:             t,
			pool:             speciesPool,
			path:             make(task.Path, 0, 2*len(t.Cities)),
			cityCount:        make([]uint32, len(t.Cities)),
			notVisitedCityID: make([]uint32, 0, len(t.Cities)),
		}
		for idx, city := range t.Cities {
			species.cityCount[city.ID] = 0
			species.notVisitedCityID[idx] = city.ID
		}

		return species
	}
	return speciesPool
}

func (pool *SpeciesPool) Get() *Species {
	return pool.Pool.Get().(*Species)
}

func (pool *SpeciesPool) Put(x *Species) {
	pool.Pool.Put(x)
}
