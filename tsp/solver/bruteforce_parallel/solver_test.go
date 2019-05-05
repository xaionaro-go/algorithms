package bruteforce

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xaionaro-go/algorithms/tsp/task"
)

func TestSolverCorrectness(t *testing.T) {
	assert.NoError(t, task.CheckSolver(New(), 8, 10*time.Second))
}

func init() {
	task.GenerateBenchmarkTasks()
}

func BenchmarkSolver_cities4(b *testing.B) {
	task.DoBenchmark(b, New(), 4)
}

func BenchmarkSolver_cities5(b *testing.B) {
	task.DoBenchmark(b, New(), 5)
}

func BenchmarkSolver_cities6(b *testing.B) {
	task.DoBenchmark(b, New(), 6)
}

func BenchmarkSolver_cities7(b *testing.B) {
	task.DoBenchmark(b, New(), 7)
}

func BenchmarkSolver_cities8(b *testing.B) {
	task.DoBenchmark(b, New(), 8)
}

func BenchmarkSolver_cities9(b *testing.B) {
	task.DoBenchmark(b, New(), 9)
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

