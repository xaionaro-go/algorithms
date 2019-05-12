package main

import (
	"testing"
)

func BenchmarkFindPermutations1(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(1)
	}
}

func BenchmarkFindPermutations2(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(2)
	}
}

func BenchmarkFindPermutations3(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(3)
	}
}

func BenchmarkFindPermutations5(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(5)
	}
}

func BenchmarkFindPermutations7(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(7)
	}
}

func BenchmarkFindPermutations10(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(10)
	}
}
/*
func BenchmarkFindPermutations15(b *testing.B) {
	for i:=0; i<b.N; i++ {
		FindPermutations(15)
	}
}
*/