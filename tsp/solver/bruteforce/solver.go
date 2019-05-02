package bruteforce

import (
	"context"
	"github.com/xaionaro-go/algorithms/tsp/task"
	"sort"
)

var _ task.Solver = New()

type Solver struct {
}

func New() *Solver {
	return &Solver{}
}

func (solver *Solver) SortTaskDataForSimple(t *task.Task) {
	for _, city := range t.Cities {
		sort.Slice(city.OutRoutes, func(i, j int) bool {
			return len(t.Cities[city.OutRoutes[i].EndCity.ID].InRoutes) < len(t.Cities[city.OutRoutes[j].EndCity.ID].InRoutes)
		})
	}
}

func (solver *Solver) SortTaskDataForFull(t *task.Task) {
	averageInRouteCost := make([]float64, len(t.Cities))
	for _, city := range t.Cities {
		sumCost := float64(0)
		for _, route := range city.InRoutes {
			sumCost += route.Cost
		}
		averageInRouteCost[city.ID] = sumCost / float64(len(city.InRoutes))
	}
	for _, city := range t.Cities {
		sort.Slice(city.OutRoutes, func(i, j int) bool {
			aRoute := city.OutRoutes[i]
			bRoute := city.OutRoutes[j]
			aCity := t.Cities[aRoute.EndCity.ID]
			bCity := t.Cities[bRoute.EndCity.ID]
			aScore := averageInRouteCost[city.ID] / aRoute.Cost / float64(len(aCity.InRoutes))
			bScore := averageInRouteCost[city.ID] / bRoute.Cost / float64(len(bCity.InRoutes))
			return aScore > bScore
		})
	}
}

/*
// It's a slow method, don't use it. Use "findSolutionSingle" instead.
func (solver *Solver) findSolutionParallel(ctxOrig context.Context, t *task.Task, parallelFactor int) task.Path {
	// The first: we need to find an any solution as fast as possible to understand some higher estimation of the cost
	solver.SortTaskDataForFull(t) // routes to "cities with less amount of incoming routes" should be first (to do less iterations on the next line)

	ctx, cancelFn := context.WithCancel(ctxOrig)
	w := newWorker(ctx, t)
	_, simplePathCost := w.findSimplePath(nil, nil, 0)

	sem := make(chan bool, parallelFactor)

	// Then we brute force all the variants with the cost lower than the estimation (from the above line)
	var path task.Path
	var result atomic.Value
	for _, divider := range []float64{1024, 128, 64, 32, 16, 8, 4, 2, 1.5, 1} {
		sem <- true
		go func(divider float64) {
			defer func() { <-sem } ()

			path, _ := w.findCheapestPath(
				nil,
				nil,
				0,
				nil,
				0,
				simplePathCost/divider,
			)
			if path != nil {
				cancelFn()
				result.Store(path)
			}
		}(divider)
		if pathI := result.Load(); pathI != nil {
			path = pathI.(task.Path)
			break
		}
	}
	if path != nil {
		return path
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	if path == nil {
		if pathI := result.Load(); pathI != nil {
			path = pathI.(task.Path)
		}
	}
	return path
}
*/

func (solver *Solver) findSolutionSingle(ctx context.Context, t *task.Task) task.Path {
	// The first: we need to find an any solution as fast as possible to understand some higher estimation of the cost
	solver.SortTaskDataForFull(t) // use the most attractive routes, first

	w := newWorker(ctx, t)
	_, simplePathCost := w.findSimplePath(nil, nil, 0)

	// Then we brute force all the variants with the cost lower than the estimation (from the above line)
	for _, divider := range []float64{1024, 128, 64, 32, 16, 8, 4, 2, 1} {
		path, _ := w.findCheapestPath(
			nil,
			nil,
			0,
			nil,
			0,
			simplePathCost/divider,
		)
		if path != nil {
			return path
		}
	}

	return nil
}

func (solver *Solver) FindSolution(ctx context.Context, t *task.Task) task.Path {
	return solver.findSolutionSingle(ctx, t)
	//return solver.findSolutionParallel(ctx, t, runtime.NumCPU())
}
