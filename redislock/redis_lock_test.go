package redislock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-pay/gopher/orm"
	"github.com/go-pay/gopher/xlog"
)

func TestClient_TryLock(t *testing.T) {
	rd := orm.InitRedis(&orm.RedisConfig{
		Addrs:    []string{"host:6379"},
		Password: "password",
		DB:       0,
	})

	locker := New(rd)
	ctx := context.Background()

	lock, err := locker.Obtain(ctx, "redis_lock_key", 20*time.Second)
	if err != nil && !errors.Is(err, ErrNotObtained) {
		xlog.Errorf("lock err: %v", err)
		return
	}
	if lock == nil {
		xlog.Warn("locker not obtained")
		return
	}
	if lock != nil {
		xlog.Warn("locker obtained")
		defer lock.Release(ctx)

		// do something

		ttl, err := lock.TTL(ctx)
		if err != nil {
			xlog.Errorf("ttl err: %v", err)
			return
		}
		xlog.Warnf("ttl: %v", ttl)

		key := lock.Key()
		xlog.Warnf("key: %s", key)

		token := lock.Token()
		xlog.Warnf("token: %s", token)
	}
}
