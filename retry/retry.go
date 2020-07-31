package retry

import (
	"time"

	"github.com/iGoogle-ink/gotil/xlog"
)

// Retry retries the callback func if some error is raised.
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
