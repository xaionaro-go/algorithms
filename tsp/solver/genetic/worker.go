package genetic

import (
	"context"
	"fmt"
	"github.com/xaionaro-go/spinlock"
	"math/rand"
	"sync/atomic"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type ShouldStopper interface {
	ShouldStop(workDone uint64, fitness float64, currentCost float64) bool
}

type worker struct {
	Ctx          context.Context
	Tick         uint32
	SpeciesPool  *SpeciesPool
	PathPool     *PathPool
	RoutesPool   *RoutesPool
	IntSlicePool *IntSlicePool
	Task         *task.Task
	LeaderLocker spinlock.Locker
	Leader       atomic.Value
}

func newWorker(ctx context.Context, t *task.Task) *worker {
	w := &worker{
		Ctx:          ctx,
		PathPool:     NewPathPool(t),
		RoutesPool:   NewRoutesPool(t),
		IntSlicePool: NewIntSlicePool(),
		Task:         t,
	}
	w.SpeciesPool = newSpeciesPool(w, t)
	w.SetLeader(w.SpeciesPool.Get())
	return w
}

func (w *worker) SetContext(newCtx context.Context) {
	w.Ctx = newCtx
}

func (w *worker) IsTimedOut() bool {
	return false
}

func (w *worker) FindRandomLoop() *task.Path {
	startCity := w.Task.Cities[rand.Intn(len(w.Task.Cities))]
	return w.FindRandomPath(startCity, startCity)

}

func (w *worker) FindRandomPath(startCity, endCity *task.City) *task.Path {
	result := w.PathPool.Get()
	cityCount := w.IntSlicePool.Get(len(w.Task.Cities), len(w.Task.Cities))

	oneMoreLoop := false
	requiredToEndASAP := w.Tick&15 == 0
	city := startCity
	possibleRoutes := w.RoutesPool.Get()
	for {
		if len(*result) > 4 {
			requiredToEndASAP = true
		}

		if requiredToEndASAP {
			for _, route := range city.OutRoutes {
				if route.EndCity == endCity {
					*result = append(*result, route)
					w.RoutesPool.Put(possibleRoutes)
					w.IntSlicePool.Put(cityCount)
					return result
				}
			}
		}

		*possibleRoutes = (*possibleRoutes)[:0]
		for _, route := range city.OutRoutes {
			if (*cityCount)[route.EndCity.ID] < 1 || (oneMoreLoop && (*cityCount)[route.EndCity.ID] < 3) {
				*possibleRoutes = append(*possibleRoutes, route)
			}
		}
		if len(*possibleRoutes) == 0 {
			requiredToEndASAP = true
			oneMoreLoop = true
			continue
		}
		r := rand.Intn(len(*possibleRoutes)) * rand.Intn(len(*possibleRoutes)) / len(*possibleRoutes)
		route := (*possibleRoutes)[r]
		*result = append(*result, route)
		if route.EndCity == endCity {
			w.RoutesPool.Put(possibleRoutes)
			w.IntSlicePool.Put(cityCount)
			return result
		}
		city = route.EndCity

		(*cityCount)[city.ID]++
	}

	panic(`should not happened`)
}

func (w *worker) SetLeader(species *Species) {
	w.Leader.Store(species)
}

func (w *worker) GetLeader() *Species {
	return w.Leader.Load().(*Species)
}

func (w *worker) considerLocalLeader(localLeader *Species) {
	w.LeaderLocker.LockDo(func() {
		currentLeader := w.GetLeader()
		if localLeader.GetFitness() > currentLeader.GetFitness() {
			fmt.Println("new leader:", localLeader.GetFitness(), "; old leader:", currentLeader.GetFitness())
			newLeader := localLeader.Clone()
			newLeader.cloneOf = nil
			w.SetLeader(newLeader)
		}
	})
}

func (w *worker) executeSubWorker(workerID int, mutateTryLimit int, newBloodFraction float64, population Population, shouldStop func(population Population, workerID int) bool) {
	var tick uint32
	for {
		tick++
		if tick&0xffff == 0 {
			select {
			case <-w.Ctx.Done(): // timeout
				return
			default:
			}
		}

		wasSuccessfulMutation := false
		localLeader := population[0]
		for idx, species := range population {
			for tryNum := 0; tryNum < mutateTryLimit; tryNum++ {
				if newSpecies := species.TryMutate(); newSpecies != nil {
					wasSuccessfulMutation = true
					population[idx] = newSpecies
					w.SpeciesPool.Put(species)
					species = newSpecies
				}
			}
			if !wasSuccessfulMutation {
				continue
			}
			if species.GetFitness() > localLeader.GetFitness() {
				localLeader = species
			}
		}
		if wasSuccessfulMutation {
			w.considerLocalLeader(localLeader)
		}
		currentLeader := w.GetLeader()
		population.CloneFrom(currentLeader, newBloodFraction)
		if shouldStop(population, workerID) {
			return
		}
	}
}

func (w *worker) Execute(parallelFactor int, shouldStopper ShouldStopper, populationSize int, mutateTryLimit int, newBloodFraction float64) {
	if parallelFactor > populationSize {
		panic(`this case is not supported, yet`)
	}
	population := NewPopulation(w.SpeciesPool, populationSize)

	epoch := make([]uint, parallelFactor)

	stopChan := make(chan struct{})
	currentEpoch := uint(0)
	shouldStop := false
	for i := 0; i < parallelFactor; i++ {
		startIdx := (len(population) * i) / parallelFactor
		endIdx := (len(population) * (i + 1)) / parallelFactor
		go w.executeSubWorker(i, mutateTryLimit, newBloodFraction, population[startIdx:endIdx], func(population Population, workerID int) bool {
			w.Tick++
			epoch[workerID]++
			epochMin := epoch[0]
			workDone := uint64(0)
			for _, workerEpoch := range epoch {
				workDone += uint64(len(population)) * uint64(workerEpoch) * uint64(mutateTryLimit)
				if workerEpoch < epochMin {
					epochMin = workerEpoch
				}
			}
			if epochMin > currentEpoch {
				currentEpoch = epochMin
				leader := w.GetLeader()
				newShouldStop := shouldStopper.ShouldStop(workDone, leader.GetFitness(), leader.path.Cost())
				if newShouldStop != shouldStop && newShouldStop == true {
					shouldStop = newShouldStop
					stopChan <- struct{}{}
				}
			}
			return shouldStop
		})
	}

	select {
	case <-w.Ctx.Done():
	case <-stopChan:
	}

	return
}
