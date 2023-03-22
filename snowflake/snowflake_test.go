package snowflake

import (
	"testing"
	"time"

	"github.com/go-pay/gopher/xlog"
)

func TestNewNode(t *testing.T) {
	xlog.Level = xlog.InfoLevel
	node, err := NewNode(1)
	if err != nil {
		xlog.Errorf("err:%v", err)
		return
	}
	for i := 0; i < 20; i++ {
		go func() {
			id := node.Generate().Int64()
			xlog.Infof("id:%d", id)
		}()
	}
	time.Sleep(time.Second)
}
