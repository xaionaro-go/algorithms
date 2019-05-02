package genetic

import (
	"math/rand"

	"github.com/xaionaro-go/algorithms/tsp/task"
)

type Species struct {
	task             *task.Task
	path             task.Path
	cityCount        []uint32
	totalCityCount   uint32
	notVisitedCityID []uint32
	isValidPath      bool
	pool             *SpeciesPool
}

func (species *Species) Fitness(t *task.Task) float64 {
	if species.isValidPath {
		// Note: positive value
		return 1 / (1 + species.path.Cost())
	}

	// Note: negative value
	return float64(len(t.Cities)-len(species.path)) * (100 - 1/(1+species.path.Cost()))
}

func (species *Species) updateIsValidPath() {
	if len(species.notVisitedCityID) != 0 {
		species.isValidPath = false
		return
	}

	species.isValidPath = true
}

func (species *Species) Clone() *Species {
	newSpecies := species.pool.Get()
	newSpecies.path = newSpecies.path[0:len(species.path)]
	copy(newSpecies.path, species.path)
	copy(newSpecies.cityCount, species.cityCount)
	newSpecies.totalCityCount = species.totalCityCount
	newSpecies.notVisitedCityID = newSpecies.notVisitedCityID[0:len(species.notVisitedCityID)]
	copy(newSpecies.notVisitedCityID, species.notVisitedCityID)
	newSpecies.isValidPath = species.isValidPath

	return newSpecies
}

func (species *Species) AddRandomRoute(routes task.Routes) {
	route := routes[rand.Intn(len(routes))]
	species.path = append(species.path, route)
	species.cityCount[route.EndCity.ID]++

	if species.cityCount[route.EndCity.ID] == 1 {
		// Remove the cityID from notVisitedCityID

		for idx, cityID := range species.notVisitedCityID {
			if cityID == route.EndCity.ID {
				species.notVisitedCityID[idx] = species.notVisitedCityID[len(species.notVisitedCityID)-1]
				species.notVisitedCityID = species.notVisitedCityID[0 : len(species.notVisitedCityID)-1]
				break
			}
		}
	}
}

func (species *Species) RandomPathPermutation(w *worker) {
	r0 := rand.Intn(len(species.path))
	r1 := rand.Intn(len(species.path) - 1)
	if r1 == r0 {
		r1++
	}
	if r1 < r0 {
		r1, r0 = r0, r1
	}
	start := species.path[r0]
	end := species.path[r1]
	newPath := w.FindRandomPath(start.StartCity, end.EndCity)

}

func (species *Species) Mutate(w *worker) {
	if len(species.path) == 0 {
		species.AddRandomRoute(w.Task.StartCity.Routes)
		return
	}

	r := rand.Intn(2 + len(species.path))

	if len(species.path) < len(w.Task.Cities) {
		if r == 0 {
			species.AddRandomRoute(w.Task.StartCity.Routes)
			return
		}
	}

	species.RandomPathPermutation(w)
}

func (species *Species) Release() {
	species.path = species.path[:0]
	species.notVisitedCityID = species.notVisitedCityID[0:len(species.cityCount)]
	for idx, city := range species.task.Cities {
		species.cityCount[city.ID] = 0
		species.notVisitedCityID[idx] = city.ID
	}
	species.isValidPath = false
	species.totalCityCount = 0
	species.pool.Put(species)
}
