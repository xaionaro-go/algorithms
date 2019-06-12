package main

import (
	"fmt"
)

func canCompleteCircuit(a []int , b []int) int {
	sumA := 0
	for _, v := range a {
		sumA += v
	}
	sumB := 0
	for _, v := range b {
		sumB += v
	}
	if sumB > sumA {
		return -1
	}

	balance := a[0]
	for idxMinusOne := range a[1:] {
		idx := idxMinusOne + 1
		balance -= b[idx-1]
	}
	if balance >= 0 {
		fmt.Println(balance)
		return 0
	}
	for idxMinusOne := range a[1:] {
		idxMinusTwo := idxMinusOne - 1
		if idxMinusTwo < 0 {
			idxMinusTwo = len(a) - 1
		}
		balance += b[idxMinusTwo] - b[idxMinusOne]
		if balance >= 0 {
			return idxMinusOne + 1
		}
	}
	return -1
}

func main() {
	fmt.Println(canCompleteCircuit([]int{1, 2}, []int{2, 1}))
}