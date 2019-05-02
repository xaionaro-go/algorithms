package genetic

import (
	"sync"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type PathPool struct {
	*sync.Pool
}

func newPathPool(t *task.Task) *PathPool {
	pathPool := &PathPool{&sync.Pool{}}
	pathPool.New = func() interface{} {
		return task.Path{}
	}
	return pathPool
}

func (pool *PathPool) Get() task.Path {
	return pool.Pool.Get().(task.Path)
}

func (pool *PathPool) Put(x task.Path) {
	x = x[:0]
	pool.Pool.Put(x)
}
