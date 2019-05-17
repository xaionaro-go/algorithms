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

	setToArray := func(idx, v int, set bool) {
		if set {
			currentArray[idx] = v
		} else {
			numCount[v]--
		}
		sign := 1
		if !set {
			sign = -1
		}
		for i := 1; i <= n; i++ {
			if i == v {
				continue
			}
			if idx+1+i < len(nonpossibleValues) {
				nIdx := idx + (1+i)*2
				if nIdx >= len(nonpossibleValues) || currentArray[nIdx] != 0 {
					nonpossibleValues[idx+1+i][i] += sign
				}
			}
			if idx-1-i >= 0 {
				nIdx := idx - (1+i)*2
				if nIdx < 0 || currentArray[nIdx] != 0 {
					nonpossibleValues[idx-1-i][i] += sign
				}
			}
		}
		if set {
			numCount[v]++
		} else {
			currentArray[idx] = 0
		}
	}

	idx := 0
	for ; currentArray[idx] != 0; idx++ {}

	for newV := 1; newV <= n; newV++ { // n
		if nonpossibleValues[idx][newV] > 0 {
			//fmt.Println("skip", newV, "not possible", nonpossibleValues[idx][newV], idx)
			continue
		}
		if numCount[newV] > 0 {
			//fmt.Println("skip", newV, "count >= 2")
			continue
		}
		setToArray(idx, newV, true)
		if idx+(newV+1) < len(currentArray) && currentArray[idx+(newV+1)] == 0 {
			setToArray(idx+(newV+1), newV, true)
			findPermutationsRecursive(n, numsLeft-2, currentArray, numCount, nonpossibleValues, results)
			setToArray(idx+(newV+1), newV, false)
		}
		if idx-(newV+1) >= 0 && currentArray[idx-(newV+1)] == 0 {
			setToArray(idx-(newV+1), newV, true)
			findPermutationsRecursive(n, numsLeft-2, currentArray, numCount, nonpossibleValues, results)
			setToArray(idx-(newV+1), newV, false)
		}
		setToArray(idx, newV, false)
	}
}

func FindPermutations(n int) (results [][]int) {
	// T: O(n!)

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
	fmt.Println(FindPermutations(9))
	fmt.Println(FindPermutations(10))
}
