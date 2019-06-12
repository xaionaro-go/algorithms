package solution

import (
	"testing"

	"github.com/xaionaro-go/algorithms/lru/task"
)

func TestCache(t *testing.T) {
	task.CheckCache(t, New())
}
