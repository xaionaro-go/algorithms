package main

import (
	"testing"
)

func doBenchmark(b *testing.B, fn func([]int, int) []int, arrayLength int, windowLength int) {
	array := make([]int, arrayLength)
	for i := 0; i < arrayLength; i++ {
		array[i] = arrayLength - i - 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(array, windowLength)
	}
}

func BenchmarkMy100000w20000(b *testing.B) {
	doBenchmark(b, slidingMaximum, 100000, 20000)
}

func BenchmarkProposed100000w20000(b *testing.B) {
	doBenchmark(b, slidingMaximum_fromSite, 100000, 20000)
}

func BenchmarkMy10000w2000(b *testing.B) {
	doBenchmark(b, slidingMaximum, 10000, 2000)
}

func BenchmarkProposed10000w2000(b *testing.B) {
	doBenchmark(b, slidingMaximum_fromSite, 10000, 2000)
}

func BenchmarkMy1000w200(b *testing.B) {
	doBenchmark(b, slidingMaximum, 1000, 200)
}

func BenchmarkProposed1000w200(b *testing.B) {
	doBenchmark(b, slidingMaximum_fromSite, 1000, 200)
}

func BenchmarkMy100w20(b *testing.B) {
	doBenchmark(b, slidingMaximum, 100, 20)
}

func BenchmarkProposed100w20(b *testing.B) {
	doBenchmark(b, slidingMaximum_fromSite, 100, 20)
}

func BenchmarkMy10w2(b *testing.B) {
	doBenchmark(b, slidingMaximum, 10, 2)
}

func BenchmarkProposed10w2(b *testing.B) {
	doBenchmark(b, slidingMaximum, 10, 2)
}
