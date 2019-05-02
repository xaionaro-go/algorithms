package genetic

import (
	"fmt"
	"math/rand"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type Species struct {
	inUse            bool
	task             *task.Task
	path             *task.Path
	cityCount        []uint32
	totalCityCount   uint32
	notVisitedCityID []uint32
	isValidPath      bool
	pool             *SpeciesPool
	cloneOf          *Species
	worker           *worker
}

func (species *Species) Reset() {
	*species.path = (*species.path)[:0]
	species.notVisitedCityID = species.notVisitedCityID[0:len(species.cityCount)]
	for idx, city := range species.task.Cities {
		species.cityCount[city.ID] = 0
		species.notVisitedCityID[idx] = city.ID
	}
	species.isValidPath = false
	species.totalCityCount = 0
	species.cloneOf = nil
}

func (species *Species) GetFitness() float64 {
	if species.isValidPath {
		// Note: positive value
		return 1 / (1 + species.path.Cost())
	}

	// Note: negative value
	return -float64(1 + len(species.task.Cities) - int(species.totalCityCount))
}

func (species *Species) updateIsValidPath() {
	if len(species.notVisitedCityID) != 0 {
		species.isValidPath = false
		return
	}

	species.isValidPath = (*species.path)[len(*species.path)-1].EndCity.ID == species.task.StartCity.ID
}

func (species *Species) Clone() *Species {
	newSpecies := species.pool.Get()
	species.CopyTo(newSpecies)
	return newSpecies
}

func (species *Species) CopyTo(dst *Species) {
	if cap(*dst.path) < len(*species.path) {
		*dst.path = make(task.Path, len(*species.path))
		if species.worker.PathPool.DefaultSize < len(*species.path) {
			species.worker.PathPool.DefaultSize = len(*species.path)
		}
	}
	*dst.path = (*dst.path)[0:len(*species.path)]
	copy(*dst.path, *species.path)
	copy(dst.cityCount, species.cityCount)
	dst.totalCityCount = species.totalCityCount
	dst.notVisitedCityID = dst.notVisitedCityID[0:len(species.notVisitedCityID)]
	copy(dst.notVisitedCityID, species.notVisitedCityID)
	dst.isValidPath = species.isValidPath
	dst.cloneOf = species
}

func (species *Species) addRandomRoute(routes task.Routes) {
	route := routes[rand.Intn(len(routes))]
	*species.path = append(*species.path, route)
	species.cityCount[route.EndCity.ID]++

	if species.cityCount[route.EndCity.ID] == 1 {
		species.totalCityCount++
		// Remove the cityID from notVisitedCityID

		for idx, cityID := range species.notVisitedCityID {
			if cityID == route.EndCity.ID {
				species.notVisitedCityID[idx] = species.notVisitedCityID[len(species.notVisitedCityID)-1]
				species.notVisitedCityID = species.notVisitedCityID[0 : len(species.notVisitedCityID)-1]
				break
			}
		}
	}
	//species.checkConsistency()
}

func (species *Species) randomPathPermutation() *task.Path {
	oldPath := species.path
	r0 := rand.Intn(len(*oldPath))
	r1 := rand.Intn(len(*oldPath))
	if r1 < r0 {
		r1, r0 = r0, r1
	}
	start := (*oldPath)[r0]
	end := (*oldPath)[r1]
	walkaroundPath := species.worker.FindRandomPath(start.StartCity, end.EndCity)
	newPath := species.worker.PathPool.Get()
	newPath.Append((*oldPath)[:r0])
	newPath.Append(*walkaroundPath)
	if len(*oldPath) > r1+1 {
		newPath.Append((*oldPath)[r1+1:])
	}
	species.path = newPath
	for _, route := range (*oldPath)[r0 : r1+1] {
		species.cityCount[route.EndCity.ID]--
		if species.cityCount[route.EndCity.ID] == 0 {
			species.notVisitedCityID = append(species.notVisitedCityID, route.EndCity.ID)
		}
	}
	for _, route := range *walkaroundPath {
		species.cityCount[route.EndCity.ID]++
		if species.cityCount[route.EndCity.ID] == 1 {
			for idx, cityID := range species.notVisitedCityID {
				if cityID == route.EndCity.ID {
					species.notVisitedCityID[idx] = species.notVisitedCityID[len(species.notVisitedCityID)-1]
					species.notVisitedCityID = species.notVisitedCityID[0 : len(species.notVisitedCityID)-1]
					break
				}
			}
		}
	}
	//species.checkConsistency()
	species.worker.PathPool.Put(walkaroundPath)
	return oldPath
}

func (species *Species) checkConsistency() {
	cityCount := make([]uint32, len(species.task.Cities))
	notVisitedCityID := make([]uint32, 0, len(species.task.Cities))

	prevCity := species.task.StartCity
	for idx, route := range *species.path {
		if route.StartCity != prevCity {
			panic(fmt.Sprint("unconsistent path: ", idx))
		}
		cityCount[route.EndCity.ID]++
		prevCity = route.EndCity
	}

	for _, city := range species.task.Cities {
		if cityCount[city.ID] == 0 {
			notVisitedCityID = append(notVisitedCityID, city.ID)
		}
	}

	for idx, _ := range cityCount {
		if cityCount[idx] != species.cityCount[idx] {
			panic(fmt.Sprint("invalid cityCount: ", idx, cityCount[idx], species.cityCount[idx], species.path))
		}
	}

	for _, cityID := range notVisitedCityID {
		found := false
		for _, cmpCityID := range species.notVisitedCityID {
			if cmpCityID == cityID {
				found = true
			}
		}
		if !found {
			panic(`invalid notVisitedCityID`)
		}
	}

	for _, cityID := range species.notVisitedCityID {
		found := false
		for _, cmpCityID := range notVisitedCityID {
			if cmpCityID == cityID {
				found = true
			}
		}
		if !found {
			panic(`invalid notVisitedCityID`)
		}
	}
}

func (species *Species) TryMutate() *Species {
	oldFitness := species.GetFitness()
	newSpecies := species.Clone()
	newSpecies.cloneOf = nil

	if len(*newSpecies.path) == 0 {
		newSpecies.addRandomRoute(species.task.StartCity.OutRoutes)
		if newSpecies.GetFitness() < oldFitness {
			fmt.Println(species.GetFitness(), oldFitness)
			// rollback
			newSpecies.Release()
			return nil
		}

		// apply
		newSpecies.updateIsValidPath()
		return newSpecies
	}

	if newSpecies.totalCityCount < uint32(len(species.task.Cities)) {
		r := rand.Intn(2 + len(*newSpecies.path))
		if r == 0 {
			newSpecies.addRandomRoute((*species.path)[len(*species.path)-1].EndCity.OutRoutes)
			if newSpecies.GetFitness() < oldFitness {
				// rollback
				newSpecies.Release()
				return nil
			}

			// apply
			newSpecies.updateIsValidPath()
			return newSpecies
		}
	}

	oldPath := newSpecies.randomPathPermutation()
	species.worker.PathPool.Put(oldPath)

	if newSpecies.GetFitness() < oldFitness {
		// rollback
		newSpecies.Release()
		return nil
	}

	// apply
	newSpecies.updateIsValidPath()
	return newSpecies
}

func (species *Species) Release() {
	species.pool.Put(species)
}
