package bruteforce

import (
	"context"
	"github.com/xaionaro-go/algorithms/tsp/task"
	"github.com/xaionaro-go/spinlock"
	"math"
	"sort"
	"time"
)

var (
	jobResultCancel = &jobResult{
		path:        nil,
		cost:        math.Inf(-1),
		ready:       true,
		dontRelease: true,
	}
	jobResultDeadEnd = &jobResult{
		path:        nil,
		cost:        -1,
		ready:       true,
		dontRelease: true,
	}
)

type workerMode int

const (
	workerModeSimplePath = workerMode(iota)
	workerModeFull
)

type worker struct {
	tick             int
	mode             workerMode
	isRunning        bool
	ctx              context.Context
	task             *task.Task
	intSlicePool     *intSlicePool
	queue            *queue
	leaderCost       AtomicFloat64
	leaderPath       task.Path
	leaderLocker     spinlock.Locker
	subWorkersCount  AtomicInt64
	parallelFactor   uint32
	jobPool          *jobPool
	jobSlicePool     *jobSlicePool
	jobArgumentsPool *jobArgumentsPool
	jobResultPool    *jobResultPool
}

func newWorker(ctx context.Context, t *task.Task) *worker {
	return &worker{
		ctx:              ctx,
		task:             t,
		intSlicePool:     newIntSlicePool(),
		jobPool:          newJobPool(),
		jobSlicePool:     newJobSlicePool(),
		jobArgumentsPool: newJobArgumentsPool(uint(len(t.Cities))),
		jobResultPool:    newJobResultPool(),
	}
}

func (w *worker) sortTaskData() {
	t := w.task

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

func (w *worker) considerSolution(subWorkerID uint32, args *jobArguments, result *jobResult) {
	//fmt.Println(subWorkerID, "solution", *args.curPath, args.curCost, w.leaderCost.GetFast())
	leaderCost := w.leaderCost.Get()
	if args.curCost >= leaderCost && leaderCost != 0 {
		return
	}
	if cap(result.path) < len(args.curPath) {
		result.path = make(task.Path, len(args.curPath))
	} else {
		result.path = result.path[:len(args.curPath)]
	}
	copy(result.path, args.curPath)
	result.cost = args.curCost
	w.leaderLocker.LockDo(func() {
		leaderCost := w.leaderCost.GetFast()
		if result.cost >= leaderCost && leaderCost != 0 {
			return
		}
		//fmt.Println(subWorkerID, "cheapest", *args.curPath, result.path, result.cost, w.leaderCost.GetFast())
		w.leaderCost.Set(result.cost)
		if cap(w.leaderPath) < len(result.path) {
			w.leaderPath = make(task.Path, len(result.path))
		} else {
			w.leaderPath = w.leaderPath[:len(result.path)]
		}
		copy(w.leaderPath, result.path)
	})
}

func (w *worker) getGlobalCostLimit() float64 {
	if !w.isRunning {
		return 0
	}
	return w.leaderCost.Get()
}

func (w *worker) enqueue(args *jobArguments, result *jobResult) *jobResult {
	job := w.jobPool.Get()
	job.args = args.Clone(w.jobArgumentsPool)
	jobResult := w.jobResultPool.Get()
	job.result = jobResult
	for !w.queue.Enqueue(job) {
		time.Sleep(time.Millisecond)
	}
	return jobResult
}

func (w *worker) stop() {
	w.isRunning = false
}

func (w *worker) newSubWorker(subWorkerID uint32) *subWorker {
	return &subWorker{
		worker: w,
		id:     subWorkerID,
	}
}

func (w *worker) Execute(parallelFactor uint32) (task.Path, bool) {
	w.isRunning = true
	w.queue = NewQueue(w.jobSlicePool, 5)
	w.parallelFactor = parallelFactor

	// The first: we need to find any solution as fast as possible to understand some higher estimation of the cost
	w.sortTaskData() // use the most attractive routes, first

	simplePathResult := w.enqueue(&jobArguments{}, nil)

	finished := make(chan struct{}, parallelFactor)
	w.mode = workerModeSimplePath
	for i := uint32(0); i < parallelFactor; i++ {
		finished <- struct{}{}
		go func() {
			defer func() { <-finished }()
			subWorker := w.newSubWorker(i)
			subWorker.execute()
		}()
	}

	doneChan := make(chan struct{})
	go func() {
		for i := uint32(0); i < parallelFactor; i++ {
			finished <- struct{}{}
		}
		doneChan <- struct{}{}
	}()

	select {
	case <-doneChan:
	case <-w.ctx.Done():
		return simplePathResult.path, false
	}

	//fmt.Println("estimation", simplePathResult.path)

	// Then we brute force all the variants with the cost lower than (or equal to) the estimation (from the above line)

	w.isRunning = true
	var results []*jobResult
	for _, divider := range []float64{1024, 128, 64, 32, 16, 8, 4, 2, 1.5, 1} { // but first we try to find a solution for my lower price (in case if the estimation was far from real)
		result := w.enqueue(&jobArguments{
			nil,
			nil,
			0,
			nil,
			0,
			simplePathResult.cost / divider,
		}, nil)
		results = append(results, result)
	}

	finished = make(chan struct{}, parallelFactor)
	w.mode = workerModeFull
	for i := uint32(0); i < parallelFactor; i++ {
		finished <- struct{}{}
		go func() {
			defer func() { <-finished }()
			subWorker := w.newSubWorker(i)
			subWorker.execute()
		}()
	}

	doneChan = make(chan struct{})
	go func() {
		for i := uint32(0); i < parallelFactor; i++ {
			finished <- struct{}{}
		}
		doneChan <- struct{}{}
	}()

	complete := false
	select {
	case <-doneChan:
		complete = true
	case <-w.ctx.Done():
		w.stop()
	}

	return w.leaderPath, complete
}
