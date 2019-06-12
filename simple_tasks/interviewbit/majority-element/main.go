package main

import (
	"fmt"
	"sort"
)

func majorityElement_fromSite(A []int )  (int) {
	len := len(A)
	major_elem, count := A[0], 1

	for i := 1; i < len; i++ {
		if count == 0 {
			major_elem,count = A[i],1
		} else if A[i] == major_elem {
			count++
		} else {
			count--
		}
	}
	return major_elem

}

func majorityElement(a []int) int {
	sort.Ints(a)
	var majorityValue *int
	count := 0
	curNum := a[0]
	for idx, v := range a {
		if v == curNum {
			count++
			continue
		}
		if count > len(a) / 2 {
			majorityValue = &a[idx-1]
		}
		count = 1
		curNum = v
	}
	if count > len(a) / 2 {
		majorityValue = &a[len(a)-1]
	}
	return *majorityValue
}

func main() {
	fmt.Println(majorityElement([]int{2, 1, 2}))
	fmt.Println(majorityElement([]int{1, 1, 2}))
	fmt.Println(majorityElement_fromSite([]int{2, 1, 2}))
	fmt.Println(majorityElement_fromSite([]int{1, 1, 2}))
}
