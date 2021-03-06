package bruteforce

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xaionaro-go/algorithms/tsp/task"
)

func TestSolverCorrectness(t *testing.T) {
	assert.NoError(t, task.CheckSolver(New(), 0, 8, 5*time.Second))
}

func init() {
	task.GenerateBenchmarkTasks()
}

func BenchmarkSolver_cities4(b *testing.B) {
	task.DoBenchmark(b, New(), 4)
}

func BenchmarkSolver_cities6(b *testing.B) {
	task.DoBenchmark(b, New(), 6)
}

func BenchmarkSolver_cities8(b *testing.B) {
	task.DoBenchmark(b, New(), 8)
}

func BenchmarkSolver_cities10(b *testing.B) {
	task.DoBenchmark(b, New(), 10)
}

func BenchmarkSolver_cities11(b *testing.B) {
	task.DoBenchmark(b, New(), 11)
}

func BenchmarkSolver_cities12(b *testing.B) {
	task.DoBenchmark(b, New(), 12)
}

func BenchmarkSolver_cities13(b *testing.B) {
	task.DoBenchmark(b, New(), 13)
}

func BenchmarkSolver_cities14(b *testing.B) {
	task.DoBenchmark(b, New(), 14)
}

func BenchmarkSolver_cities15(b *testing.B) {
	task.DoBenchmark(b, New(), 15)
}
