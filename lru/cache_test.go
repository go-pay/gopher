package lru

import (
	"testing"

	"github.com/go-pay/gopher/xlog"
)

func TestNewCache(t *testing.T) {
	xlog.Level = xlog.InfoLevel
	cache := NewCache(2)

	cache.Put("1", "one")
	xlog.Info(cache.Get("1"))

	xlog.Warn("===============")
	cache.Put("2", "two")
	xlog.Info(cache.Get("1"))
	xlog.Info(cache.Get("2"))

	xlog.Warn("===============")
	cache.Put("3", "three")
	xlog.Info(cache.Get("1")) // nil
	xlog.Info(cache.Get("2"))
	xlog.Info(cache.Get("3"))

	xlog.Warn("===============")
	cache.Put("1", "one")
	xlog.Info(cache.Get("2")) // nil
	xlog.Info(cache.Get("3"))
	xlog.Info(cache.Get("1"))

	xlog.Warn("===============")
	cache.Put("2", "two")
	xlog.Info(cache.Get("1"))
	xlog.Info(cache.Get("2"))
	xlog.Info(cache.Get("3")) // nil
}
