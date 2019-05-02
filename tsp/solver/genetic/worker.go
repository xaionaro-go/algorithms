package genetic

import (
	"context"
	"math/rand"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type worker struct {
	SpeciesPool *SpeciesPool
	PathPool    *PathPool
	RoutesPool  *RoutesPool
	Task        *task.Task
	Leader      *Species
}

func newWorker(t *task.Task) *worker {
	speciesPool := newSpeciesPool(t)
	leader := speciesPool.Get()

	return &worker{
		SpeciesPool: speciesPool,
		PathPool:    newPathPool(t),
		RoutesPool:  newRoutesPool(t),
		Task:        t,
		Leader:      leader,
	}
}

func (w *worker) FindRandomLoop() task.Path {
	startCity := w.Task.Cities[rand.Intn(len(w.Task.Cities))]
	return w.FindRandomPath(startCity, startCity)

}

func (w *worker) FindRandomPath(startCity, endCity *task.City) task.Path {
	result := w.PathPool.Get()

	city := startCity
	for {
		for _, route := range city.Routes {
			if route.EndCity == endCity {
				result = append(result, route)
				return result
			}
		}

		possibleRoutes := w.RoutesPool.Get()
		route := city.Routes[rand.Intn(len(city.Routes))]

		for _, cmpRoute := range result {
			if route == cmpRoute {
				alreadyHasTheRoute = true
			}
		}
		result = append(result, route)
		city = route.EndCity
	}

	panic(`should not happened`)
}

func (w *worker) Execute(ctx context.Context) {

}
