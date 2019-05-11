package task

import (
	"context"
	"testing"
)

const (
	MaxBenchmarkCityAmount = 20
)

var (
	benchmarkTasks [1001][]*Task
)

func GenerateBenchmarkTasks() {
	for i := 3; i <= MaxBenchmarkCityAmount; i++ {
		for _, seed := range []int64{0, 1, 2, 3, 4, 5, 6, 7} {
			task, _ := GenerateRiddle(uint32(i), seed)
			benchmarkTasks[i] = append(benchmarkTasks[i], task)
		}
	}
	for _, i := range []int{25, 30, 35, 40, 45, 50, 60, 70, 80, 90, 95, 100, 120, 250, 1000} {
		for _, seed := range []int64{0, 1, 2, 3, 4, 5, 6, 7} {
			task, _ := GenerateRiddle(uint32(i), seed)
			benchmarkTasks[i] = append(benchmarkTasks[i], task)
		}
	}
}

func DoBenchmark(b *testing.B, solver Solver, cityAmount int) {
	if len(benchmarkTasks[cityAmount]) == 0 {
		panic(`there're no pre-generated tasks for this amount if cities'`)
	}
	b.ReportAllocs()
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, task := range benchmarkTasks[cityAmount] {
			solver.FindSolution(ctx, task)
		}
	}
}
