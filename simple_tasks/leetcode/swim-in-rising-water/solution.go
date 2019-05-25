package main

import (
	"container/heap"
)

type point struct {
	x int
	y int
	t int
}

func (p *point) EqualsTo(cmp *point) bool {
	return p.x == cmp.x && p.y == cmp.y
}

func (p *point) IsValidOnGrid(grid [][]int) bool {
	if p.x < 0 || p.y < 0 {
		return false
	}

	if p.y > len(grid)-1 {
		return false
	}

	if p.x > len(grid[p.y])-1 {
		return false
	}

	return true
}

type pointsQueue []*point

func (q *pointsQueue) Len() int {
	return len(*q)
}

func (q *pointsQueue) Swap(i, j int) {
	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

func (q *pointsQueue) Less(i, j int) bool {
	return (*q)[i].t < (*q)[j].t
}

func (q *pointsQueue) Push(p interface{}) {
	*q = append(*q, p.(*point))
}

func (q *pointsQueue) Pop() interface{} {
	p := (*q)[q.Len()-1]
	*q = (*q)[:q.Len()-1]
	return p
}

func (q *pointsQueue) PushPoint(p *point) {
	heap.Push(q, p)
}

func (q *pointsQueue) PopPoint() *point {
	p := heap.Pop(q)
	if p == nil {
		return nil
	}

	return p.(*point)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func swimInWater(grid [][]int) int {
	if len(grid) == 0 {
		return -1
	}
	if len(grid[0]) == 0 {
		return -1
	}
	endPos := &point{
		x: len(grid[len(grid)-1]) - 1,
		y: len(grid) - 1,
	}
	alreadyWasGrid := make([][]bool, len(grid))
	for idx, line := range grid {
		alreadyWasGrid[idx] = make([]bool, len(line))
	}
	var queue pointsQueue
	queue.PushPoint(&point{0, 0, grid[0][0]})
	for {
		pos := queue.PopPoint()
		if pos == nil {
			break
		}
		alreadyWasGrid[pos.x][pos.y] = true
		for _, diff := range []point{{-1, 0, 0}, {1, 0, 0}, {0, -1, 0}, {0, 1, 0}} {
			newPos := &point{
				x: pos.x + diff.x,
				y: pos.y + diff.y,
			}
			if !newPos.IsValidOnGrid(grid) {
				continue
			}
			if alreadyWasGrid[newPos.x][newPos.y] {
				continue
			}
			newPos.t = max(grid[newPos.y][newPos.x], pos.t)
			if newPos.EqualsTo(endPos) {
				return newPos.t
			}
			queue.PushPoint(newPos)
		}
	}
	panic(`Should not happened`)
}
