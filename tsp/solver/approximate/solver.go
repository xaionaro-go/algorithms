package approximate

import (
	"context"
	"github.com/xaionaro-go/algorithms/tsp/task"
	"math"
	"sort"
)

var _ task.Solver = New()

type Solver struct {
}

func New() *Solver {
	return &Solver{}
}

func (solver *Solver) SortTaskData(t *task.Task) {
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
	solver.SortTaskData(t)
	w := newWorker(ctx, t)

	requiredCityCount := make([]int, len(t.Cities))
	for i := 0; i < len(t.Cities); i++ {
		requiredCityCount[i] = 1
	}

	path, _ := w.findPath(
		t.StartCity,
		t.StartCity,
		requiredCityCount,
	)

	if path == nil {
		return nil
	}

	path = w.optimizePath(path)
	return path
}
