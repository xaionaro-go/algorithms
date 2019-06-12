package main

import (
	"fmt"
)

func firstStepNum(base int) int {
	firstDigit := base
	order := 1
	for firstDigit >= 10 {
		firstDigit /= 10
		order *= 10
	}

	result := 0
	prevDigit := firstDigit
	order /= 10
	for order > 0 {
		curDigit := base / order % 10
		if curDigit != prevDigit - 1 && curDigit != prevDigit + 1 {
			if curDigit < prevDigit + 1 {
				if curDigit < prevDigit - 1 {
					curDigit = prevDigit - 1
				} else {
					curDigit = prevDigit + 1
				}
			} else {
				
			}
			for order > 0 {
				result += curDigit * order
				order /= 10
				curDigit--
				if curDigit < 0 {
					curDigit = 1
				}
			}
			return result
		}
		result += base / order * order
		order /= 10
	}

	return nextStepNum(base)
}

func nextStepNum(num int) int {
	// 508
	// 543
	// 545
	// 565
	// 567
	// 654
	// 656
	// 676
	// 678
	// ...
	// 1010
	// 1012
	// 1210
	// 1212

	nextStepNum := num % 10
	prevDigit := num % 10

	k := 10
	for {
		num /= 10
		digit := num % 10
		if digit < prevDigit && prevDigit < 9 {
			digit = prevDigit+1

		}
		nextStepNum += digit * k
		k *= 10
		prevDigit = digit
	}

	return nextStepNum
}

func stepnum(a, b int) (result []int) {
	num := firstStepNum(a)
	for num <= b {
		result = append(result, num)
		num = nextStepNum(a)
	}
	return
}

func main() {
	fmt.Println(stepnum(508, 1100))
}