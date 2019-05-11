package approximate

import (
	"github.com/stretchr/testify/assert"
	"github.com/xaionaro-go/algorithms/tsp/task"
	"testing"
	"time"
)

func TestSolverCorrectness(t *testing.T) {
	assert.NoError(t, task.CheckSolver(New(), 0.0, 4, 5*time.Second))
	assert.NoError(t, task.CheckSolver(New(), 0.0, 8, 5*time.Second))
	assert.NoError(t, task.CheckSolver(New(), 0.05, 12, 5*time.Second))
	assert.NoError(t, task.CheckSolver(New(), 0.05, 20, 5*time.Second))
	assert.NoError(t, task.CheckSolver(New(), 0.0, 25, 5*time.Second))
	assert.NoError(t, task.CheckSolver(New(), 0.0, 30, 5*time.Second))
}

func init() {
	task.GenerateBenchmarkTasks()
}

func BenchmarkSolver_cities4(b *testing.B) {
	task.DoBenchmark(b, New(), 4)
}

/*
func BenchmarkSolver_cities6(b *testing.B) {
	task.DoBenchmark(b, New(), 6)
}
*/
func BenchmarkSolver_cities8(b *testing.B) {
	task.DoBenchmark(b, New(), 8)
}

func BenchmarkSolver_cities10(b *testing.B) {
	task.DoBenchmark(b, New(), 10)
}

/*
func BenchmarkSolver_cities11(b *testing.B) {
	task.DoBenchmark(b, New(), 11)
}
*/
func BenchmarkSolver_cities12(b *testing.B) {
	task.DoBenchmark(b, New(), 12)
}

/*
func BenchmarkSolver_cities13(b *testing.B) {
	task.DoBenchmark(b, New(), 13)
}

func BenchmarkSolver_cities14(b *testing.B) {
	task.DoBenchmark(b, New(), 14)
}
*/
func BenchmarkSolver_cities15(b *testing.B) {
	task.DoBenchmark(b, New(), 15)
}

func BenchmarkSolver_cities20(b *testing.B) {
	task.DoBenchmark(b, New(), 20)
}

func BenchmarkSolver_cities25(b *testing.B) {
	task.DoBenchmark(b, New(), 25)
}

func BenchmarkSolver_cities30(b *testing.B) {
	task.DoBenchmark(b, New(), 30)
}

/*
func BenchmarkSolver_cities35(b *testing.B) {
	task.DoBenchmark(b, New(), 35)
}
*/
func BenchmarkSolver_cities40(b *testing.B) {
	task.DoBenchmark(b, New(), 40)
}

/*
func BenchmarkSolver_cities45(b *testing.B) {
	task.DoBenchmark(b, New(), 45)
}
*/
func BenchmarkSolver_cities50(b *testing.B) {
	task.DoBenchmark(b, New(), 50)
}

func BenchmarkSolver_cities60(b *testing.B) {
	task.DoBenchmark(b, New(), 60)
}

func BenchmarkSolver_cities70(b *testing.B) {
	task.DoBenchmark(b, New(), 70)
}

func BenchmarkSolver_cities80(b *testing.B) {
	task.DoBenchmark(b, New(), 80)
}

func BenchmarkSolver_cities90(b *testing.B) {
	task.DoBenchmark(b, New(), 90)
}

func BenchmarkSolver_cities95(b *testing.B) {
	task.DoBenchmark(b, New(), 95)
}

/*
func BenchmarkSolver_cities100(b *testing.B) {
	task.DoBenchmark(b, New(), 100)
}

func BenchmarkSolver_cities120(b *testing.B) {
	task.DoBenchmark(b, New(), 120)
}
*/
/*
func BenchmarkSolver_cities250(b *testing.B) {
	task.DoBenchmark(b, New(), 250)
}

func BenchmarkSolver_cities1000(b *testing.B) {
	task.DoBenchmark(b, New(), 1000)
}
*/
