package approximate

import (
	"sync"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type pathPool struct {
	*sync.Pool
	DefaultSize uint
}

func newPathPool() *pathPool {
	pathPool := &pathPool{Pool: &sync.Pool{}}
	pathPool.New = func() interface{} {
		return &[]task.Path{{}}[0]
	}
	return pathPool
}

func (pool *pathPool) Get(l uint) *task.Path {
	r := pool.Pool.Get().(*task.Path)
	if l == 0 {
		l = pool.DefaultSize
	}
	if uint(cap(*r)) < l {
		r = &[]task.Path{make(task.Path, 0, l)}[0]
	}
	return r
}

func (pool *pathPool) Put(x *task.Path) {
	*x = (*x)[:0]
	pool.Pool.Put(x)
}
