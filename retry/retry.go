package retry

import (
	"time"

	"github.com/go-pay/gopher/xlog"
)

// Retry 重试 func 最大次数，间隔
func Retry(callback func() error, maxRetries int, interval time.Duration) (err error) {
	for i := 1; i <= maxRetries; i++ {
		if err = callback(); err != nil {
			xlog.Warnf("Retry(%d) err(%+v)", i, err)
			time.Sleep(interval)
			continue
		}
		return
	}
	return
}
