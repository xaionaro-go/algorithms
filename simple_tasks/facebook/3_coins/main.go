package main

import (
	"fmt"
	"sort"
)

func FindAllCoinSums_dynamic(coinTypes []int, upToSum int, fn func(sum int)) {
	// T: O(t*s + s*log(s))
	// M: O(s)

	hasSumMap := map[int]bool{}
	hasSumMap[0] = true

	for sum := 1; sum <= upToSum; sum++ {
		for _, coinType := range coinTypes {
			prevSum := sum - coinType
			hasSum := hasSumMap[prevSum]
			if hasSum {
				hasSumMap[sum] = true
			}
		}
	}

	result := make([]int, 0, len(hasSumMap))
	for k := range hasSumMap {
		result = append(result, k)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	for _, sum := range result {
		if sum == 0 {
			continue
		}
		fn(sum)
	}
}

func findAllCoinSums_recursive(coinTypes []int, bottomLimit *int, upToSum int, fn func(sum int), curSum int) {
	*bottomLimit = curSum
	for _, coinType := range coinTypes {
		newSum := curSum + coinType
		if newSum > upToSum {
			continue
		}
		if newSum <= *bottomLimit {
			continue
		}
		panic(`wrong order`)
		fn(newSum)
	}
	for _, coinType := range coinTypes {
		newSum := curSum + coinType
		if newSum > upToSum {
			continue
		}
		if newSum <= *bottomLimit {
			continue
		}
		findAllCoinSums_recursive(coinTypes, bottomLimit, upToSum, fn, curSum+coinType)
	}
}

func FindAllCoinSums_recursive(coinTypes []int, upToSum int, fn func(sum int)) {
	// Works incorrectly
	//
	// T: O(s + t*log(t))
	// M: O(s)

	sort.Ints(coinTypes)

	findAllCoinSums_recursive(coinTypes, &[]int{0}[0], upToSum, fn, 0)
}

func FindAllCoinSums_increment(coinTypes []int, upToSum int, fn func(sum int)) {
	// Works incorrectly, it's required to fix coin conversion
	//
	// T: O(t^4 * s)
	// M: O(t^2)

	hasSumMap := map[int]bool{}

	sort.Ints(coinTypes)
	coinCount := make([]int, len(coinTypes))
	coinCountTransaction := make([]int, len(coinTypes))
	bestCoinCountTransaction := make([]int, len(coinTypes))
	sum := 0
	for sum < upToSum {
		prevSum := sum
		nextSum := upToSum + 1
		for coinIdx, coinType := range coinTypes {
			for idx := range coinCountTransaction {
				coinCountTransaction[idx] = 0
			}

			sum = prevSum + coinType
			coinCountTransaction[coinIdx]++
			//fmt.Println("A", prevSum, sum, nextSum, coinCount, coinCountTransaction)

			if sum < nextSum {
				copy(bestCoinCountTransaction, coinCountTransaction)
				nextSum = sum
			}

			for {
				sumBeforeSub := nextSum
				//fmt.Println("S", sumBeforeSub, nextSum, coinCount)
				for coinIdxSub, coinTypeSub := range coinTypes[:coinIdx] {
					if coinIdxSub == coinIdx {
						continue
					}
					if coinCount[coinIdxSub] + coinCountTransaction[coinIdxSub] <= 0 {
						if !hasSumMap[sum-coinIdxSub-coinType] {
							continue
						}
					}
					sum -= coinTypeSub
					if sum <= prevSum {
						sum += coinTypeSub
						break
					}
					if sum >= nextSum {
						sum += coinTypeSub
						continue
					}
					coinCountTransaction[coinIdxSub]--
					copy(bestCoinCountTransaction, coinCountTransaction)
					nextSum = sum
				}
				//fmt.Println("R", sumBeforeSub, nextSum, coinCount)
				if nextSum == sumBeforeSub {
					break
				}
			}

		}
		if nextSum > upToSum {
			break
		}

		// This's incorrect way to convert coins (we should remove common dividers)
		for coinIdx, coinType := range coinTypes {
			if coinCount[coinIdx] > 10 {
				continue
			}
			for coinIdxConv, coinTypeConv := range coinTypes {
				if coinIdxConv == coinIdx {
					continue
				}
				if coinCount[coinIdxConv]+bestCoinCountTransaction[coinIdxConv] > coinType {
					bestCoinCountTransaction[coinIdxConv] -= coinType
					bestCoinCountTransaction[coinIdx] += coinTypeConv
				}
			}
		}

		notEnoughOfSomeCoin := false
		for idx, count := range bestCoinCountTransaction {
			coinCount[idx] += count
			if coinCount[idx] <= 1 {
				notEnoughOfSomeCoin = true
			}
		}
		fn(nextSum)
		if notEnoughOfSomeCoin {
			hasSumMap[nextSum] = true
		}
		sum = nextSum
	}
}

func main() {
	FindAllCoinSums_increment([]int{10, 15, 55}, 1000, func(sum int) {
		fmt.Println(sum)
	})
}
