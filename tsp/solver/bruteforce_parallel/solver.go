package bruteforce

import (
	"context"
	"runtime"

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
	path, complete := w.Execute(uint32(runtime.NumCPU() / 2))
	if !complete {
		return nil
	}
	return path
}

func (solver *Solver) FindBestSolutionBeforeCancel(ctx context.Context, t *task.Task) (task.Path, bool) {
	w := newWorker(ctx, t)
	return w.Execute(uint32(runtime.NumCPU() / 2))
}
