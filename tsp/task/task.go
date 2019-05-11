package task

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xaionaro-go/errors"
)

const (
	MaxBranching = 100
)

type City struct {
	ID        uint32
	OutRoutes Routes
	InRoutes  Routes
}
type Cities []*City

func (s Cities) Sort() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID < s[j].ID
	})
}

type Path []*Route

func (path Path) Cost() (result float64) {
	if path == nil {
		return 0
	}
	for _, route := range path {
		result += route.Cost
	}
	return
}

func (path Path) String() string {
	var nodes []string
	if len(path) > 0 {
		nodes = append(nodes, strconv.FormatInt(int64(path[0].StartCity.ID), 10))
	}
	for _, route := range path {
		nodes = append(nodes, strconv.FormatInt(int64(route.EndCity.ID), 10))
	}

	return strings.Join(nodes, `>`) + " " + strconv.FormatFloat(path.Cost(), 'f', 5, 64)
}

func (path *Path) Append(appendPath Path) *Path {
	*path = append(*path, appendPath...)
	return path
}

type Route struct {
	StartCity *City
	EndCity   *City
	Cost      float64
}

type Routes []*Route

func (s Routes) Shuffle() {
	perm := rand.Perm(len(s))
	for i, v := range perm {
		s[i], s[v] = s[v], s[i]
	}
}

type Task struct {
	StartCity *City
	Cities    Cities
}

func (t *Task) IsValidPath(path Path) bool {
	if len(path) < len(t.Cities) {
		return false
	}

	if path[0].StartCity != t.StartCity {
		return false
	}

	if path[len(path)-1].EndCity != t.StartCity {
		return false
	}

	return true
}

type Solver interface {
	FindSolution(ctx context.Context, task *Task) Path
}

// GenerateRiddle returns a task (startCity, routes) and the solution for the task.
func GenerateRiddle(cityAmount uint32, seed int64) (task *Task, solution Path) {
	task = &Task{}

	randGen := rand.New(rand.NewSource(seed))

	cityIDs := randGen.Perm(int(cityAmount))

	for _, cityID := range cityIDs {
		task.Cities = append(task.Cities, &City{
			ID: uint32(cityID),
		})
	}

	for cityIdx, city := range task.Cities {
		/*prevCityIdx := cityIdx
		if prevCityIdx == 0 {
			prevCityIdx = len(cityIDs)-1
		} else {
			prevCityIdx--
		}
		prevCity := task.Cities[prevCityIdx]*/

		nextCityIdx := cityIdx + 1
		if nextCityIdx >= len(task.Cities) {
			nextCityIdx = 0
		}

		nextCity := task.Cities[nextCityIdx]

		// How many routes there will be from this city
		branching := uint32(randGen.Intn(MaxBranching))

		// We shouldn't produce more routes from a city than the count of other cities
		if branching > cityAmount-1 {
			branching = cityAmount - 1
		}

		for i := uint32(0); i < branching; i++ {
			endCityIdx := randGen.Intn(int(cityAmount) - 2)
			if endCityIdx >= cityIdx {
				endCityIdx += 2 // skip current city and the next one
			}
			endCityIdx %= len(task.Cities)

			endCity := task.Cities[endCityIdx]
			route := &Route{
				StartCity: city,
				EndCity:   endCity,
				Cost:      rand.Float64(),
			}
			// multiple routes from between two cities is OK (there could be multiple flights from a city to another one)
			city.OutRoutes = append(city.OutRoutes, route)
			endCity.InRoutes = append(city.InRoutes, route)
		}

		// Generating a route for the winner-path (solution)
		route := &Route{
			StartCity: city,
			EndCity:   nextCity,
			Cost:      rand.Float64() * rand.Float64(), // Not a strict way to win, but I hope it will work on big amount of cities
		}
		city.OutRoutes = append(city.OutRoutes, route)
		nextCity.InRoutes = append(nextCity.InRoutes, route)
		solution = append(solution, route)
	}

	for _, city := range task.Cities {
		city.OutRoutes.Shuffle()
		city.InRoutes.Shuffle()
	}

	task.StartCity = task.Cities[0]
	task.Cities.Sort()

	return
}

func checkSolution(task *Task, expectedSolution Path, proposedSolution Path, accuracyBacklash float64) errors.SmartError {
	if len(proposedSolution) < len(task.Cities) {
		return errors.InvalidArguments.Wrap(`invalid length of the solution`,
			fmt.Sprint(len(proposedSolution), ` != `, len(expectedSolution)), proposedSolution, expectedSolution)
	}

	if expectedSolution.Cost() < proposedSolution.Cost()*(1 - accuracyBacklash) {
		return errors.New(`bad solution (to high cost)`, expectedSolution, proposedSolution)
	}

	if !task.IsValidPath(proposedSolution) {
		return errors.New(`invalid path`, proposedSolution, task.StartCity.ID, len(task.Cities))
	}

	return nil
}

func CheckSolver(solver Solver, accuracyBacklash float64, cityAmount uint32, durationLimit time.Duration) error {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(durationLimit))
	for _, seed := range []int64{0, 1, 2, 3, 4, 5, 6, 7} {
		task, solution := GenerateRiddle(cityAmount, seed)
		err := checkSolution(task, solution, solver.FindSolution(ctx, task), accuracyBacklash)
		if err != nil {
			return err.Wrap(`seed:`, seed)
		}
	}
	return nil
}

// CheckTimeComplexity returns:
// * NaN -- if the execution was just too long
// * 0 -- if constant complexity
// * 1 -- if linear complexity
// * 2 -- if square complexity
// * if log(n) then the returned value will be >0 and <1
// * and so on
//
// If the value is greater than "4" then it's most probably an exponential complexity
func CheckTimeComplexity(solver Solver, durationLimit time.Duration) (result float64) {
	result = math.NaN()
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(durationLimit))

	inputSizes := []int{3, 4, 5, 10, 30, 100, 300, 1000, 3000, 10000}
	var execDuration, prevExecDuration time.Duration
	for idx, inputSize := range inputSizes {
		prevExecDuration = execDuration

		task, _ := GenerateRiddle(uint32(inputSize), 0)
		startTime := time.Now()
		solver.FindSolution(ctx, task)

		select {
		case <-ctx.Done():
			return // timeout
		default:
		}

		execDuration = time.Since(startTime)

		if idx == 0 {
			continue
		}
		prevInputSize := inputSizes[idx-1]

		result = math.Log(float64(execDuration.Nanoseconds())/float64(prevExecDuration.Nanoseconds())) /
			math.Log(float64(inputSize)/float64(prevInputSize))
	}

	return
}
