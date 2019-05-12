package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNode_Add(t *testing.T) {
	for _, testdata := range []struct {
		i0 uint
		i1 uint
	}{
		{
			i0: 1234,
			i1: 1234,
		},
		{
			i0: 997979,
			i1: 12345,
		},
	} {
		list0 := GenerateList(testdata.i0)
		list1 := GenerateList(testdata.i1)

		assert.Equal(t, int(testdata.i0+testdata.i1), list0.Add(list1).ToInt())
	}
}

func benchmarkNode_Add(b *testing.B, pow int) {
	list0 := GenerateList(uint(math.Pow10(pow) + 0.5))
	list1 := GenerateList(uint(math.Pow10(pow) + 0.5))
	list1.digit = 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list0.Add(list1)
	}
}

func BenchmarkNode_Add_10e5(b *testing.B) {
	benchmarkNode_Add(b, 5)
}

func BenchmarkNode_Add_10e10(b *testing.B) {
	benchmarkNode_Add(b, 10)
}

func BenchmarkNode_Add_10e15(b *testing.B) {
	benchmarkNode_Add(b, 15)
}
