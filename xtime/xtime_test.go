package xtime

import (
	"testing"
	"time"

	"github.com/iGoogle-ink/gopher/xlog"
)

func TestXtime(t *testing.T) {
	minutes := Time(1609066441).Time().Add(time.Minute * 30).Sub(time.Now()).Minutes()
	xlog.Debug(minutes)
	if minutes < 0 { // 30分钟超时
		//更新订单状态为订单超时
		xlog.Debug("超时")
	}
}
