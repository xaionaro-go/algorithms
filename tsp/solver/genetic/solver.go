package genetic

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sort"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

var _ task.Solver = New()

type Solver struct {
}

func New() *Solver {
	return &Solver{}
}

type stopper struct {
	NoChangeLimit uint
	noChangeCount uint
	previousCost  float64
}

func (s *stopper) ShouldStop(workDone uint64, fitness float64, cost float64) bool {
	fmt.Println(s.noChangeCount, workDone, fitness, cost)
	if cost == s.previousCost {
		s.noChangeCount++
	} else {
		s.previousCost = cost
		s.noChangeCount = 0
	}
	return s.noChangeCount > s.NoChangeLimit
}

func (solver *Solver) SortTaskDataForFull(t *task.Task) {
	minimalInRouteCost := make([]float64, len(t.Cities))
	for _, city := range t.Cities {
		min := city.InRoutes[0].Cost
		for _, route := range city.InRoutes {
			if min > route.Cost {
				min = route.Cost
			}
		}
		minimalInRouteCost[city.ID] = min
	}
	for _, city := range t.Cities {
		sort.Slice(city.OutRoutes, func(i, j int) bool {
			aRoute := city.OutRoutes[i]
			bRoute := city.OutRoutes[j]
			aCity := t.Cities[aRoute.EndCity.ID]
			bCity := t.Cities[bRoute.EndCity.ID]
			aScore := (1 + minimalInRouteCost[aCity.ID]/aRoute.Cost) / aRoute.Cost / aRoute.Cost / (1 + math.Log(float64(len(aCity.InRoutes))))
			bScore := (1 + minimalInRouteCost[bCity.ID]/bRoute.Cost) / bRoute.Cost / bRoute.Cost / (1 + math.Log(float64(len(bCity.InRoutes))))
			return aScore > bScore
		})
	}
}

func (solver *Solver) FindSolution(ctx context.Context, t *task.Task) task.Path {
	solver.SortTaskDataForFull(t)
	worker := newWorker(ctx, t)
	worker.Execute(runtime.NumCPU(), &stopper{NoChangeLimit: uint(len(t.Cities)) / 2}, 64*len(t.Cities)*len(t.Cities), 2*len(t.Cities), 0.5)
	leader := worker.GetLeader()
	path := *leader.path
	fmt.Println(path)
	if t.IsValidPath(path) {
		return path
	}
	if leader.isValidPath {
		leader.updateIsValidPath()
		panic(leader.isValidPath)
	}
	return nil
}
