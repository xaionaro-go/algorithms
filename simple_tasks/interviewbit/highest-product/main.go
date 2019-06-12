package main

import (
	"fmt"
	"sort"
)

func maxp3(a []int) int {
	sort.Slice(a, func(i,j int) bool {
		i0 := a[i]
		i1 := a[j]
		if i0 < 0 {
			i0 = -i0
		}
		if i1 < 0 {
			i1 = -i1
		}
		return i0 > i1
	})
	res := a[0] * a[1] * a[2]
	if res >= 0 {
		return res
	}
	max := -1
	for _, v := range a[3:] {
		for _, cmpValue := range[]int{a[0] * a[1] * v, a[0] * a[2] * v, a[1] * a[2] * v} {
			if cmpValue > max {
				max = cmpValue
			}
		}
	}
	if max >= 0 {
		return max
	}
	return a[len(a)-3] * a[len(a)-2] * a[len(a)-1]
}

func main() {
	fmt.Println(maxp3([]int{0, -1, 3, -100, -49, 70, 50}))
	fmt.Println(maxp3([]int{58, -31, -24, 58, 25, -31, 28, 2, -21, -48}))
}
