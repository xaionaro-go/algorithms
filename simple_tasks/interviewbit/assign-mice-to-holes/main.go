package main

import (
	"fmt"
	"sort"
)

func mice(a []int, b []int) int {
	sort.Ints(a)
	sort.Ints(b)
	maxDiff := 0
	for idx := range a {
		diff := b[idx] - a[idx]
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff
}

func main() {
	fmt.Println(mice([]int{4, -4, 2}, []int{4, 0, 5}))
}
