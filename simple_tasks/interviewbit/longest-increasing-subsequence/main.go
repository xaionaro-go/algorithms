package main

import (
	"fmt"
)


func lis(a []int) int {
	var m = map[int]int{}
	var recursive func(a []int, curIdx int) int
	recursive = func(a []int, curIdx int) int {
		if len(a) <= curIdx + 1 {
			return 0
		}

		if res, ok := m[curIdx]; ok {
			return res
		}

		max := 0
		for idx, v := range a[curIdx+1:] {
			if curIdx >= 0 && v <= a[curIdx] {
				continue
			}

			var res int
			if idx+1 < len(a) {
				res = 1 + recursive(a, curIdx+1+idx)
			} else {
				res = 1
			}
			if res > max {
				max = res
			}
		}
		m[curIdx] = max
		return max
	}
	return recursive(a, -1)
}

func main() {
	fmt.Println(lis([]int{1, 3, 5}))
	fmt.Println(lis([]int{0, 8, 4, 12, 2, 10, 6, 14, 1, 9, 5, 13, 3, 11, 7, 15}))
}
