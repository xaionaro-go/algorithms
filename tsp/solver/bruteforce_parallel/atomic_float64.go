package bruteforce

import (
	"math"
	"sync/atomic"
	"unsafe"
)

// AtomicFloat64 is an implementation of atomic float64 using uint64 atomic instructions
// and `math.Float64frombits()`/`math.Float64bits()`
type AtomicFloat64 float64

// Get returns the current value
func (f *AtomicFloat64) Get() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)((unsafe.Pointer)(f))))
}

// Set sets a new value
func (f *AtomicFloat64) Set(n float64) {
	atomic.StoreUint64((*uint64)((unsafe.Pointer)(f)), math.Float64bits(n))
}

// Add adds the value to the current one (operator "plus")
func (f *AtomicFloat64) Add(a float64) float64 {
	for {
		// Get the old value
		o := f.Get()

		// Calculate the sum
		s := o + a

		// Get int64 representation of the sum to be able to use atomic operations
		n := math.Float64bits(s)

		// Swap the old value to the new one
		// If not successful then somebody changes the value while our calculations above
		// It means we need to recalculate the new value and try again (that's why it's in the loop)
		if atomic.CompareAndSwapUint64((*uint64)((unsafe.Pointer)(f)), math.Float64bits(o), n) {
			return s
		}
	}
}

// Add adds the value to the current one (operator "plus")
func (f *AtomicFloat64) SetIfX(n float64, fn func(oldValue float64) bool) float64 {
	for {
		// Get the old value
		o := f.Get()

		if !fn(o) {
			return o
		}

		// Swap the old value to the new one
		// If not successful then it means that meanwhile somebody changes the value
		// It means we need to recalculate the new value and try again (that's why it's in the loop)
		if atomic.CompareAndSwapUint64((*uint64)((unsafe.Pointer)(f)), math.Float64bits(o), math.Float64bits(n)) {
			return n
		}
	}
}

func (f *AtomicFloat64) SetIfLess(n float64) float64 {
	return f.SetIfX(n, func(oldValue float64) bool {
		return n < oldValue
	})
}

// GetFast is like Get but without atomicity (faster, but unsafe)
func (f *AtomicFloat64) GetFast() float64 {
	return float64(*f)
}

// SetFast is like Set but without atomicity (faster, but unsafe)
func (f *AtomicFloat64) SetFast(n float64) {
	*f = AtomicFloat64(n)
}

// AddFast is like Add but without atomicity (faster, but unsafe)
func (f *AtomicFloat64) AddFast(n float64) float64 {
	*f += AtomicFloat64(n)
	return float64(*f)
}
