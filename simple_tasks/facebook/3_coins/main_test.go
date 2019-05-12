package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncrement(t *testing.T) {
	totalSumDynamic := 0
	totalCountDynamic := 0
	FindAllCoinSums_dynamic([]int{10, 15, 55}, 100000, func(sum int) {
		totalSumDynamic += sum
		totalCountDynamic++
	})

	totalSumIncrement := 0
	totalCountIncrement := 0
	FindAllCoinSums_increment([]int{10, 15, 55}, 100000, func(sum int) {
		totalSumIncrement += sum
		totalCountIncrement++
	})

	assert.Equal(t, totalSumDynamic, totalSumIncrement)
	assert.Equal(t, totalCountDynamic, totalCountIncrement)
}

func BenchmarkFindAllCoinSums_dynamic_100000(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindAllCoinSums_dynamic([]int{10, 15, 55}, 100000, func(sum int) {})
	}
}

func BenchmarkFindAllCoinSums_increment_100000(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindAllCoinSums_increment([]int{10, 15, 55}, 100000, func(sum int) {})
	}
}

func BenchmarkFindAllCoinSums_dynamic_1000000(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindAllCoinSums_dynamic([]int{10, 15, 55}, 1000000, func(sum int) {})
	}
}

func BenchmarkFindAllCoinSums_increment_1000000(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindAllCoinSums_increment([]int{10, 15, 55}, 1000000, func(sum int) {})
	}
}

func BenchmarkFindAllCoinSums_dynamic_10000000(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindAllCoinSums_dynamic([]int{10, 15, 55}, 10000000, func(sum int) {})
	}
}

func BenchmarkFindAllCoinSums_increment_10000000(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindAllCoinSums_increment([]int{10, 15, 55}, 10000000, func(sum int) {})
	}
}