package main

import "fmt"

import (
	"container/heap"
)

type maxHeap struct {
	Items []*maxHeapItem
}

func newMaxHeap(windowLength int) *maxHeap {
	return &maxHeap{
		Items: make([]*maxHeapItem, 0, windowLength*2+1),
	}
}

type maxHeapItem struct {
	Index     int
	HeapIndex int
	Value     int
}

func (heap *maxHeap) Less(i, j int) bool {
	if heap.Items[i] == nil && heap.Items[j] == nil {
		return false
	}
	if heap.Items[i] == nil {
		return false
	}
	if heap.Items[j] == nil {
		return true
	}
	return heap.Items[i].Value > heap.Items[j].Value
}

func (heap *maxHeap) Swap(i, j int) {
	heap.Items[i].HeapIndex, heap.Items[j].HeapIndex = heap.Items[j].HeapIndex, heap.Items[i].HeapIndex
	heap.Items[i], heap.Items[j] = heap.Items[j], heap.Items[i]
}

func (heap *maxHeap) Len() int {
	return len(heap.Items)
}

func (heap *maxHeap) Push(itemI interface{}) {
	item := itemI.(*maxHeapItem)
	item.HeapIndex = len(heap.Items)
	heap.Items = append(heap.Items, item)
}

func (heap *maxHeap) Pop() interface{} {
	r := heap.Items[len(heap.Items)-1]
	heap.Items = heap.Items[:len(heap.Items)-1]
	return r
}

type solutionWorker struct {
	windowLength    int
	maxHeap         *maxHeap
	heapItemByIndex []*maxHeapItem
}

func newSolutionWorker(arrayLength, windowLength int) *solutionWorker {
	return &solutionWorker{
		windowLength:    windowLength,
		maxHeap:         newMaxHeap(windowLength),
		heapItemByIndex: make([]*maxHeapItem, arrayLength),
	}
}

func (worker *solutionWorker) addToHeap(index, value int) {
	/*fmt.Print(index, " AB: ")
	for _, item := range worker.maxHeap.Items {
		fmt.Print(item, " ")
	}
	fmt.Println()*/

	item := &maxHeapItem{Index: index, Value: value}
	heap.Push(worker.maxHeap, item)
	worker.heapItemByIndex[index] = item

	/*fmt.Print(index, " AR: ")
	for _, item := range worker.maxHeap.Items {
		fmt.Print(item, " ")
	}
	fmt.Println()*/
}

/*
func (worker *solutionWorker) removeFromHeapByIndex(index int) {
	heap := worker.maxHeap
	removeItem := worker.heapItemByIndex[index]
	heapIndex := removeItem.HeapIndex
	for heapIndex < worker.windowLength {
		if heapIndex*2+1 >= len(heap.Items) {
			heap.Items[heapIndex] = nil
			break
		}
		highestChild := heap.Items[heapIndex*2+1]
		if highestChild == nil || (highestChild != nil && len(heap.Items) > heapIndex*2+2 && heap.Items[heapIndex*2+2] != nil && heap.Items[heapIndex*2+2].Value > highestChild.Value) {
			highestChild = heap.Items[heapIndex*2+2]
		}
		heap.Items[heapIndex] = highestChild
		if highestChild == nil {
			break
		}
		nextIndex := highestChild.HeapIndex
		highestChild.HeapIndex = heapIndex
		heapIndex = nextIndex
	}

	heap.Items[heapIndex] = heap.Items[len(heap.Items)-1]
	heap.Items = heap.Items[:len(heap.Items)-1]
	if heapIndex >= len(heap.Items) {
		return
	}
	for heapIndex > 0 {
		child := heap.Items[heapIndex]
		heapIndex /= 2
		parent := heap.Items[heapIndex]
		if child.Value <= parent.Value {
			break
		}
		heap.Items[parent.HeapIndex], heap.Items[child.HeapIndex] = heap.Items[child.HeapIndex], heap.Items[parent.HeapIndex]
		parent.HeapIndex, child.HeapIndex = child.HeapIndex, parent.HeapIndex
	}
}
*/

func (worker *solutionWorker) removeFromHeapByIndex(index int) {
	heap := worker.maxHeap
	/*fmt.Print(index, " B: ")
	for _, item := range heap.Items {
		fmt.Print(item, " ")
	}
	fmt.Println()*/

	removeItem := worker.heapItemByIndex[index]
	lastItem := heap.Items[len(heap.Items)-1]
	heap.Items[lastItem.HeapIndex] = nil
	lastItem.HeapIndex = removeItem.HeapIndex
	heap.Items[lastItem.HeapIndex] = lastItem
	heap.Items = heap.Items[:len(heap.Items)-1]

	{
		heapIndex := lastItem.HeapIndex
		for heapIndex > 0 && heapIndex < len(heap.Items) {
			parent := heap.Items[heapIndex/2]
			if lastItem.Value > parent.Value {
				heap.Swap(lastItem.HeapIndex, parent.HeapIndex)
			}
			heapIndex /= 2
		}
	}

	heapIndex := lastItem.HeapIndex
	parent := lastItem
	for heapIndex*2+1 < len(heap.Items) {
		child0 := heap.Items[heapIndex*2+1]
		maxChild := child0
		if len(heap.Items) > heapIndex*2+2 {
			child1 := heap.Items[heapIndex*2+2]
			if child1.Value > maxChild.Value {
				maxChild = child1
			}
		}
		if maxChild.Value <= parent.Value {
			break
		}
		heapIndex = maxChild.HeapIndex
		maxChild.HeapIndex, parent.HeapIndex = parent.HeapIndex, maxChild.HeapIndex
		heap.Items[maxChild.HeapIndex], heap.Items[parent.HeapIndex] = heap.Items[parent.HeapIndex], heap.Items[maxChild.HeapIndex]
		parent = heap.Items[heapIndex]
	}
	/*fmt.Print(index, " R: ")
	for _, item := range heap.Items {
		fmt.Print(item, " ")
	}
	fmt.Println()*/
}

func (worker *solutionWorker) Solve(a []int) []int {
	result := make([]int, len(a)-worker.windowLength+1)
	for idx := 0; idx < worker.windowLength; idx++ {
		worker.addToHeap(idx, a[idx])
	}
	result[0] = worker.maxHeap.Items[0].Value
	for idx := range a[:len(a)-worker.windowLength] {
		addIdx := idx + worker.windowLength
		removeIdx := idx
		worker.removeFromHeapByIndex(removeIdx)
		worker.addToHeap(addIdx, a[addIdx])
		result[idx+1] = worker.maxHeap.Items[0].Value
	}
	return result
}

func slidingMaximum(A []int, B int) []int {
	worker := newSolutionWorker(len(A), B)
	return worker.Solve(A)
}

func slidingMaximum_fromSite(A []int, B int) []int {

	max := 0
	start := true

	for i := 0; i < B; i++ {
		if start || A[i] > max {
			max = A[i]
			start = false
		}
	}
	maxes := []int{max}

	for ienter := B; ienter < len(A); ienter++ {
		iexit := ienter - B

		if A[ienter] >= max {
			max = A[ienter]
		} else if A[iexit] == max {
			start := true
			for k := iexit + 1; k <= ienter; k++ {
				if start || A[k] > max {
					max = A[k]
					start = false
				}
			}
		}
		maxes = append(maxes, max)
	}

	return maxes
}

func main() {
	fmt.Println(slidingMaximum([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3))
	fmt.Println(slidingMaximum([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, 2))
	fmt.Println(slidingMaximum([]int{90, 943, 777, 658, 742, 559, 623, 263, 880, 176, 354, 434, 699, 501, 551, 821, 563, 974, 701, 479, 238, 87, 61, 910, 204, 534, 369, 845, 566, 19, 939, 87, 708, 323, 662, 32, 655, 835, 67, 360, 550, 173, 488, 420, 680, 805, 630, 48, 791, 991, 791, 819, 772, 228, 123, 303, 642, 780, 115, 89, 919, 830, 271, 853, 588, 249, 20, 940, 851, 749, 340, 587, 235, 106, 125, 32, 319, 590, 354, 751, 761, 564, 484, 51, 202, 370, 216, 130, 146, 632}, 6))
	fmt.Println(slidingMaximum([]int{263, 215, 169, 328, 32, 735, 521, 3, 337, 705, 272, 509, 865, 60, 268, 604, 707, 803, 692, 603, 405, 994, 472, 702, 468, 841, 457, 421, 285, 90, 844, 957, 866, 220, 71, 731, 454, 204, 146, 97, 907, 727, 321, 836, 525, 877, 795, 491, 754, 238, 709, 355, 293, 326, 99, 640, 427, 93, 743, 137}, 6))
	fmt.Println(slidingMaximum_fromSite([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3))
	fmt.Println(slidingMaximum_fromSite([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, 2))
	fmt.Println(slidingMaximum_fromSite([]int{90, 943, 777, 658, 742, 559, 623, 263, 880, 176, 354, 434, 699, 501, 551, 821, 563, 974, 701, 479, 238, 87, 61, 910, 204, 534, 369, 845, 566, 19, 939, 87, 708, 323, 662, 32, 655, 835, 67, 360, 550, 173, 488, 420, 680, 805, 630, 48, 791, 991, 791, 819, 772, 228, 123, 303, 642, 780, 115, 89, 919, 830, 271, 853, 588, 249, 20, 940, 851, 749, 340, 587, 235, 106, 125, 32, 319, 590, 354, 751, 761, 564, 484, 51, 202, 370, 216, 130, 146, 632}, 6))
	fmt.Println(slidingMaximum_fromSite([]int{263, 215, 169, 328, 32, 735, 521, 3, 337, 705, 272, 509, 865, 60, 268, 604, 707, 803, 692, 603, 405, 994, 472, 702, 468, 841, 457, 421, 285, 90, 844, 957, 866, 220, 71, 731, 454, 204, 146, 97, 907, 727, 321, 836, 525, 877, 795, 491, 754, 238, 709, 355, 293, 326, 99, 640, 427, 93, 743, 137}, 6))
}
