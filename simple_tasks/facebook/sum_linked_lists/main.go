package main

import (
	"fmt"
)

type node struct {
	digit uint8
	next  *node
}

func (root *node) ToInt() (result int) {
	for n := root; n != nil; n = n.next {
		result *= 10
		result += int(n.digit)
	}
	return
}

func (root *node) reverse() *node {
	start := root
	prev := start
	next := start
	for next != nil {
		nextnext := next.next
		next.next = start
		start = next
		prev.next = nextnext
		next = nextnext
	}

	return start
}

func (root0 *node) Add(root1 *node) *node {
	// T: O(num_digits)
	// M: O(1)

	root0 = root0.reverse()
	root1 = root1.reverse()

	n0 := root0
	n1 := root1

	for n0.next != nil && n1 != nil {
		n0.digit += n1.digit
		if n0.digit >= 10 {
			n0.digit -= 10
			if n0.next == nil {
				n0.next = &node{}
			}
			n0.next.digit++
		}
		n0 = n0.next
		n1 = n1.next
	}

	for n1 != nil {
		n0.digit += n1.digit
		if n0.digit >= 10 {
			n0.digit -= 10
			n0.next = &node{}
			n0.next.digit++
		}
		n0 = n0.next
		n1 = n1.next
		if n1 != nil && n0.next == nil {
			n0.next = &node{}
		}
	}

	root1 = root1.reverse()
	root0 = root0.reverse()
	return root0
}

func GenerateList(in uint) (result *node) {
	for in > 0 {
		prevResult := result
		result = &node{
			digit: uint8(in % 10),
			next:  prevResult,
		}
		in /= 10
	}
	return result
}

func main() {
	list0 := GenerateList(997979)
	list1 := GenerateList(12345)

	fmt.Println(list0.Add(list1), 997979+12345)
}
