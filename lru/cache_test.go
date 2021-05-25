package lru

import (
	"testing"

	"github.com/iGoogle-ink/gotil/xlog"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(2)

	cache.Put("1", "one")
	xlog.Info(cache.Get("1"))

	xlog.Warn("===============")
	cache.Put("2", "two")
	xlog.Info(cache.Get("1"))

	xlog.Warn("===============")
	cache.Put("3", "three")
	xlog.Info(cache.Get("2"))
	xlog.Info(cache.Get("3"))
	xlog.Info(cache.Get("3"))
	xlog.Info(cache.Get("1"))

	xlog.Warn("===============")
	cache.Put("2", "two")
	xlog.Info(cache.Get("3"))
	xlog.Info(cache.Get("1"))
}
