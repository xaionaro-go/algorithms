package genetic

import (
	"sync"
)

type IntSlicePool struct {
	*sync.Pool
}

func NewIntSlicePool() *IntSlicePool {
	return &IntSlicePool{
		Pool: &sync.Pool{
			New: func() interface{} {
				return &[][]int{[]int(nil)}[0]
			},
		},
	}
}

func (pool *IntSlicePool) Get(l, c int) *[]int {
	s := pool.Pool.Get().(*[]int)
	if cap(*s) < c {
		*s = make([]int, c)
	}
	if len(*s) != l {
		*s = (*s)[0:l]
	}
	return s
}

func (pool *IntSlicePool) Put(s *[]int) {
	(*s) = (*s)[:cap(*s)]
	for idx := range *s {
		(*s)[idx] = 0
	}
	pool.Pool.Put(s)
}
