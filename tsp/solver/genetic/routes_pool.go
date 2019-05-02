package genetic

import (
	"sync"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type RoutesPool struct {
	*sync.Pool
}

func newRoutesPool(t *task.Task) *RoutesPool {
	routesPool := &RoutesPool{&sync.Pool{}}
	routesPool.New = func() interface{} {
		return task.Routes{}
	}
	return routesPool
}

func (pool *RoutesPool) Get() task.Routes {
	return pool.Pool.Get().(task.Routes)
}

func (pool *RoutesPool) Put(x task.Routes) {
	x = x[:0]
	pool.Pool.Put(x)
}
