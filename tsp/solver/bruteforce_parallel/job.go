package bruteforce

import (
	"github.com/xaionaro-go/algorithms/tsp/task"
)

type jobArguments struct {
	cityCount        []int
	uselessCityCount []int
	totalCityCount   int
	curPath          task.Path
	curCost          float64
	costLimit        float64
}

func (args *jobArguments) Release(pool *jobArgumentsPool) {
	args.Reset()
	pool.Put(args)
}

func (args *jobArguments) Reset() {
	return
	if args.cityCount != nil {
		args.cityCount = args.cityCount[:0]
	}
	if args.uselessCityCount != nil {
		args.uselessCityCount = args.uselessCityCount[:0]
	}
	args.totalCityCount = 0
	if args.curPath != nil {
		args.curPath = args.curPath[:0]
	}
	args.curCost = 0
	args.costLimit = 0
}

func (args *jobArguments) Clone(argsPool *jobArgumentsPool) *jobArguments {
	c := argsPool.Get()
	c.totalCityCount = args.totalCityCount
	c.curCost = args.curCost
	c.costLimit = args.costLimit

	if args.cityCount != nil {
		c.cityCount = c.cityCount[0:len(args.cityCount)]
		copy(c.cityCount, args.cityCount)
	}

	if args.uselessCityCount != nil {
		c.uselessCityCount = c.uselessCityCount[0:len(args.uselessCityCount)]
		copy(c.uselessCityCount, args.uselessCityCount)
	}

	if args.curPath != nil {
		if cap(c.curPath) < len(args.curPath) {
			c.curPath = make(task.Path, len(args.curPath), cap(args.curPath))
		} else {
			c.curPath = c.curPath[0:len(args.curPath)]
		}
		copy(c.curPath, args.curPath)
	}
	return c
}

type jobResult struct {
	path        task.Path
	cost        float64
	ready       bool
	dontRelease bool
}

func (result *jobResult) PointTo(src *jobResult) {
	result.path = src.path
	result.cost = src.cost
}

func (result *jobResult) CopyFrom(src *jobResult) {
	result.cost = src.cost

	if cap(result.path) < len(src.path) {
		result.path = make(task.Path, len(src.path))
	} else {
		result.path = result.path[:len(src.path)]
	}

	copy(result.path, src.path)
}

func (result *jobResult) Release(pool *jobResultPool) {
	pool.Put(result)
}

type job struct {
	inUse  bool
	args   *jobArguments
	result *jobResult
}

func (j *job) Release(pool *jobPool) {
	pool.Put(j)
}
