package retry

import (
	"time"

	"github.com/iGoogle-ink/gotil/xlog"
)

// Retry 重试 func 最大次数，间隔
func Retry(callback func() error, maxRetries int, interval time.Duration) error {
	var err error
	for i := 1; i <= maxRetries; i++ {
		if err = callback(); err != nil {
			xlog.Warningf("Retry(%d) error(%+v)", i, err)
			time.Sleep(interval)
			continue
		}
		return nil
	}
	return err
}
