package main

import (
	"fmt"
)

func trailingZeroes(A int)  (int) {
	if A < 5 {
		return 0
	}

	countOf5 := 0

	pow5 := 1
	for pow5 <= A {
		pow5 *= 5
		countOf5 += A / pow5
	}

	return countOf5
}


func main() {
	fmt.Println(trailingZeroes(1))
	fmt.Println(trailingZeroes(9))
	fmt.Println(trailingZeroes(10))
	fmt.Println(trailingZeroes(24))
	fmt.Println(trailingZeroes(25))
	fmt.Println(trailingZeroes(30))
	fmt.Println(trailingZeroes(1000))
}