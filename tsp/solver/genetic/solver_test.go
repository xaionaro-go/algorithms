package genetic

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xaionaro-go/algorithms/tsp/task"
)

func TestSolverTimeComplexity(t *testing.T) {
	timeComplexity := task.CheckTimeComplexity(New(), time.Second)
	assert.True(t, timeComplexity < -1, fmt.Sprintf("%v", timeComplexity))
}

func BenchmarkSolver_cities4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task.CheckSolver(New(), 4, time.Second)
	}
}

func BenchmarkSolver_cities6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task.CheckSolver(New(), 6, time.Second)
	}
}

func BenchmarkSolver_cities8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task.CheckSolver(New(), 8, time.Second)
	}
}
