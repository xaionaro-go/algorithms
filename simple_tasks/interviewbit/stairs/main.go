package main

import (
	"fmt"
)

var m = map[int]int{}
func stairsHowManyWays_recursive(stepsCount int) int {
	if stepsCount == 0 {
		return 1
	}
	if stepsCount < 0 {
		return 0
	}
	res, ok := m[stepsCount]
	if ok {
		return res
	}
	res = stairsHowManyWays_recursive(stepsCount - 1) + stairsHowManyWays_recursive(stepsCount - 2)
	m[stepsCount] = res
	return res
}


func main() {
	fmt.Println(stairsHowManyWays_recursive(3))
	fmt.Println(stairsHowManyWays_recursive(10))
}