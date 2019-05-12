package main

import (
	"testing"
)

func BenchmarkFindPermutations1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(1)
	}
}

func BenchmarkFindPermutations2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(2)
	}
}

func BenchmarkFindPermutations3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(3)
	}
}

func BenchmarkFindPermutations4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(4)
	}
}

func BenchmarkFindPermutations5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(5)
	}
}

func BenchmarkFindPermutations6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(6)
	}
}

func BenchmarkFindPermutations7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(7)
	}
}

func BenchmarkFindPermutations8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(8)
	}
}

func BenchmarkFindPermutations9(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindPermutations(9)
	}
}

func BenchmarkFindPermutations10(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(10)
	}
}

func BenchmarkFindPermutations11(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(11)
	}
}

func BenchmarkFindPermutations12(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(12)
	}
}

/*
func BenchmarkFindPermutations13(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(13)
	}
}
*/