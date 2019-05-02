package genetic

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
	worker := newWorker(t)
	worker.Execute(ctx)
	return worker.Leader.path
}
