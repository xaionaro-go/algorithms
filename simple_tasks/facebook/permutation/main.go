package main

import (
	"fmt"
)

func findPermutationsRecursive(n int, numsLeft int, currentArray []int, numCount []int, nonpossibleValues [][]int, results *[][]int) {
	//fmt.Println(n, numsLeft, currentArray, numCount, nonpossibleValues, *results)

	if numsLeft == 0 {
		arrayCopy := make([]int, len(currentArray))
		copy(arrayCopy, currentArray)
		*results = append(*results, arrayCopy)
		return
	}

	idx := n*2 - numsLeft
	for newV := 1; newV <= n; newV++ { // n
		if nonpossibleValues[idx][newV] > 0 {
			continue
		}
		if numCount[newV] >= 2 {
			continue
		}
		if numCount[newV] == 1 {
			hasThePair := false
			if idx-newV-1 >= 0 {
				if currentArray[idx-newV-1] == newV {
					hasThePair = true
				}
			}
			if idx+newV+1 < len(currentArray) {
				if currentArray[idx+newV+1] == newV {
					hasThePair = true
				}
			}
			if !hasThePair {
				continue
			}
		}
		currentArray[idx] = newV
		for i := 1; i <= n; i++ {
			if i == newV {
				continue
			}
			if idx+1+i < len(nonpossibleValues) {
				nonpossibleValues[idx+1+i][i]++
			}
			if idx-1-i >= 0 {
				nonpossibleValues[idx-1-i][i]++
			}
		}
		numCount[newV]++
		findPermutationsRecursive(n, numsLeft-1, currentArray, numCount, nonpossibleValues, results) // depth: O(n)
		numCount[newV]--
		currentArray[idx] = 0

		for i := 1; i <= n; i++ {
			if i == newV {
				continue
			}
			if idx+1+i < len(nonpossibleValues) {
				nonpossibleValues[idx+1+i][i]--
			}
			if idx-1-i >= 0 {
				nonpossibleValues[idx-1-i][i]--
			}
		}
	}
}

func FindPermutations(n int) (results [][]int) {
	nonpossibleValues := make([][]int, n*2)
	for idx := range nonpossibleValues {
		nonpossibleValues[idx] = make([]int, n+1)
	}
	numCount := make([]int, n+1)
	findPermutationsRecursive(n, n*2, make([]int, n*2), numCount, nonpossibleValues, &results)
	return
}

func main() {
	fmt.Println(FindPermutations(2))
	fmt.Println(FindPermutations(3))
	fmt.Println(FindPermutations(4))
	fmt.Println(FindPermutations(5))
	fmt.Println(FindPermutations(6))
	fmt.Println(FindPermutations(7))
	fmt.Println(FindPermutations(8))
}
