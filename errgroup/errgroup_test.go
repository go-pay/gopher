package errgroup

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-pay/xlog"
)

func TestErrgroup(t *testing.T) {
	var count int64 = 1
	countBackup := count
	eg := WithContext(context.Background())

	// go 协程
	eg.Go(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	})
	// go 协程
	eg.Go(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	})
	// go 协程
	eg.Go(func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return errors.New("error ,reset count")
	})
	// wait 协程 Done
	if err := eg.Wait(); err != nil {
		// do some thing
		count = countBackup
		xlog.Error(err)
		//return
	}
	xlog.Debug(count)
}

func TestErrgroupLimit1(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	var (
		eg    Group
		goNum = 3 // every times run goNum goroutine
	)
	for i := 0; i < 11; i++ {
		var count = int64(i)
		eg.Go(func(ctx context.Context) error {
			atomic.AddInt64(&count, 1)
			xlog.Debug("count:", count)
			return nil
		})
		if eg.WorkNum() == goNum {
			if err := eg.Wait(); err != nil {
				xlog.Error("err1:", err)
				// to do something you need
			}
			xlog.Info("wait")
			time.Sleep(time.Second)
		}
	}
	if err := eg.Wait(); err != nil {
		xlog.Error("err2:", err)
	}
}

func TestErrgroupLimit2(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	var (
		eg Group
	)
	eg.GOMAXPROCS(3)
	for i := 0; i < 11; i++ {
		var count = int64(i)
		eg.Go(func(ctx context.Context) error {
			atomic.AddInt64(&count, 1)
			xlog.Debug("count:", count)
			time.Sleep(time.Second * 3)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		xlog.Error("err2:", err)
	}
}

func timeSleep1(c context.Context) error {
	data := make(chan string, 1)
	go func() {
		time.Sleep(1 * time.Second)
		data <- "timeSleep1 done"
	}()
	select {
	case <-c.Done():
		xlog.Infof("timeSleep1 cancel")
	case rsp := <-data:
		xlog.Info(rsp)
	}
	return nil
}

func timeErrSleep3(c context.Context) error {
	time.Sleep(3 * time.Second)
	return errors.New("timeSleep3 error")
}

func timeTimeoutSleep3(c context.Context) error {
	data := make(chan string, 1)
	go func() {
		time.Sleep(3 * time.Second)
		data <- "timeSleep3 done"
	}()
	select {
	case <-c.Done():
		xlog.Infof("timeSleep3 timeout")
	case rsp := <-data:
		xlog.Info(rsp)
	}
	return nil
}

func timeCancelSleep5(c context.Context) error {
	data := make(chan string, 1)
	go func() {
		time.Sleep(5 * time.Second)
		data <- "timeCancelSleep5 done"
	}()
	select {
	case <-c.Done():
		xlog.Infof("timeCancelSleep5 cancel")
	case rsp := <-data:
		xlog.Info(rsp)
	}
	return nil
}

func timeTimeoutSleep5(c context.Context) error {
	data := make(chan string, 1)
	go func() {
		time.Sleep(5 * time.Second)
		data <- "timeTimeoutSleep5 done"
	}()
	select {
	case <-c.Done():
		return errors.New("timeTimeoutSleep5 timeout")
	case rsp := <-data:
		xlog.Info(rsp)
	}
	return nil
}

func TestErrgroupWithCancel(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	var (
		eg = WithCancel(context.Background())
	)
	eg.Go(func(ctx context.Context) error {
		err := timeSleep1(ctx)
		return err
	})
	eg.Go(func(ctx context.Context) error {
		if err := timeErrSleep3(ctx); err != nil {
			xlog.Error(err)
			return err
		}
		return nil
	})
	eg.Go(func(ctx context.Context) error {
		if err := timeCancelSleep5(ctx); err != nil {
			xlog.Errorf("err:%v", err)
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		// do some thing
	}
}

func TestErrgroupWithTimeout(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	var (
		eg = WithTimeout(context.Background(), time.Second*4)
	)
	eg.Go(func(ctx context.Context) error {
		err := timeSleep1(ctx)
		return err
	})
	eg.Go(func(ctx context.Context) error {
		err := timeTimeoutSleep3(ctx)
		return err
	})
	eg.Go(func(ctx context.Context) error {
		if err := timeTimeoutSleep5(ctx); err != nil {
			xlog.Errorf("err:%v", err)
			return err
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		// do some thing
	}
}
