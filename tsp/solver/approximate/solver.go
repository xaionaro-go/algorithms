package approximate

import (
	"context"
	"github.com/xaionaro-go/algorithms/tsp/task"
)

var _ task.Solver = New()

type Solver struct {
}

func New() *Solver {
	return &Solver{}
}

func (solver *Solver) FindSolution(ctx context.Context, t *task.Task) task.Path {
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
