package bruteforce

import (
	"sync/atomic"
)

type AtomicInt64 int64

// Get returns the current value
func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

// Set sets a new value
func (i *AtomicInt64) Set(n int64) {
	atomic.StoreInt64((*int64)(i), n)
}

// Add adds the value to the current one (operator "plus")
func (i *AtomicInt64) Add(a int64) int64 {
	return atomic.AddInt64((*int64)(i), a)
}
