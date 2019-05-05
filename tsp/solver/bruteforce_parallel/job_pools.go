package bruteforce

import (
	"sync"
)

const (
	memoryReuse = true
)

type jobArgumentsPool struct {
	cityAmount uint
	pool       *sync.Pool
}

func newJobArgumentsPool(cityAmount uint) *jobArgumentsPool {
	pool := &jobArgumentsPool{
		cityAmount: cityAmount,
		pool:       &sync.Pool{},
	}

	pool.pool.New = func() interface{} {
		return &jobArguments{
			cityCount:        make([]int, pool.cityAmount),
			uselessCityCount: make([]int, pool.cityAmount),
		}
	}
	return pool
}

func (pool *jobArgumentsPool) Get() *jobArguments {
	return pool.pool.Get().(*jobArguments)
}

func (pool *jobArgumentsPool) Put(x *jobArguments) {
	if !memoryReuse {
		return
	}
	pool.pool.Put(x)
}

type jobResultPool struct {
	pool *sync.Pool
}

func newJobResultPool() *jobResultPool {
	return &jobResultPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return &jobResult{}
			},
		},
	}
}

func (pool *jobResultPool) Get() *jobResult {
	r := pool.pool.Get().(*jobResult)
	//fmt.Printf("%p get()\n", r)
	return r
}

func (pool *jobResultPool) Put(x *jobResult) {
	if !memoryReuse {
		return
	}
	if x.dontRelease {
		return
	}
	//fmt.Printf("%p put()\n", x)
	if !x.ready {
		panic(`should not happened`)
	}
	x.ready = false
	x.path = x.path[:0]
	pool.pool.Put(x)
}

type jobPool struct {
	//nextJobID uint64
	pool *sync.Pool
}

func newJobPool() *jobPool {
	r := &jobPool{
		pool: &sync.Pool{},
	}
	r.pool.New = func() interface{} {
		return &job{}
	}
	return r
}

func (pool *jobPool) Get() *job {
	j := pool.pool.Get().(*job)
	if j.inUse {
		panic(`should not happened`)
	}
	j.inUse = true
	//j.id = atomic.AddUint64(&pool.nextJobID, 1)
	return j
}

func (pool *jobPool) Put(x *job) {
	if !memoryReuse {
		return
	}
	if !x.inUse {
		panic(`should not happened`)
	}
	x.inUse = false
	x.args = nil
	x.result = nil
	pool.pool.Put(x)
}

type jobSlicePool struct {
	pool *sync.Pool
}

func newJobSlicePool() *jobSlicePool {
	return &jobSlicePool{
		pool: &sync.Pool{
			New: func() interface{} {
				return []*job{}
			},
		},
	}
}

func (pool *jobSlicePool) Get(size uint) []*job {
	s := pool.pool.Get().([]*job)
	if uint(cap(s)) < size {
		s = make([]*job, size)
		s = s[:0]
	}
	return s
}

func (pool *jobSlicePool) Put(x []*job) {
	if !memoryReuse {
		return
	}
	x = x[:0]
	pool.pool.Put(x)
}
