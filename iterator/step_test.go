package iterator

import (
	"testing"

	"github.com/go-pay/gopher/xlog"
)

func TestSteps(t *testing.T) {
	var ids = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for idx, step := range Steps(len(ids), 6) {
		batchIds := ids[step.Head:step.Tail]
		xlog.Warnf("[%d]step, slice: %d", idx, batchIds)
	}
}
