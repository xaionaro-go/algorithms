package genetic

import (
	"sync"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type SpeciesPool struct {
	*sync.Pool
	worker *worker
}

func newSpeciesPool(w *worker, t *task.Task) *SpeciesPool {
	speciesPool := &SpeciesPool{Pool: &sync.Pool{}, worker: w}
	speciesPool.New = func() interface{} {
		species := &Species{
			task:             t,
			worker:           speciesPool.worker,
			pool:             speciesPool,
			path:             &[]task.Path{make(task.Path, 0, 2*len(t.Cities))}[0],
			cityCount:        make([]uint32, len(t.Cities)),
			notVisitedCityID: make([]uint32, len(t.Cities)),
		}
		species.Reset()

		return species
	}
	return speciesPool
}

func (pool *SpeciesPool) Get() *Species {
	species := pool.Pool.Get().(*Species)
	if species.inUse {
		panic(`Concurrent use of one species :(`)
	}
	species.inUse = true
	return species
}

func (pool *SpeciesPool) Put(x *Species) {
	x.inUse = false
	pool.Pool.Put(x)
}
