package main

import (
	"fmt"
)

func main() {
	fmt.Println(mincostToHireWorkers([]int{10, 20, 5}, []int{70, 50, 30}, 2))
	fmt.Println(mincostToHireWorkers([]int{3, 1, 10, 10, 1}, []int{4, 8, 2, 2, 7}, 3))
}
