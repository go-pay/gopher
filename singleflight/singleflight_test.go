package singleflight

import (
	"context"
	"testing"
	"time"

	"github.com/go-pay/gopher/orm"
	"github.com/go-pay/xlog"
)

type Single struct {
	Name string
	Age  int
}

func TestDo(t *testing.T) {
	var g Group[*Single]
	// 启动10个协程去请求数据，最终执行的只有1个，输出一次 I am return Single
	for i := 0; i < 10; i++ {
		go func() {
			v, _, err := g.Do("key", func() (*Single, error) {
				xlog.Warn("working produce Single")
				s := &Single{
					Name: "lady",
					Age:  18,
				}
				return s, nil
			})
			if err != nil {
				t.Errorf("Do error = %v", err)
			}
			xlog.Warnf("final result: %#v", v)
		}()
	}
	time.Sleep(2 * time.Second)
}

func TestDoChan(t *testing.T) {
	var g Group[*Single]
	// 启动10个协程去请求数据，最终执行的只有1个，输出一次 I am return Single
	for i := 0; i < 10; i++ {
		go func() {
			doChan := g.DoChan("key", func() (*Single, error) {
				xlog.Warn("working produce Single")
				s := &Single{
					Name: "lady",
					Age:  18,
				}
				return s, nil
			})
			result := <-doChan
			xlog.Warnf("final result: %#v", result.Val)
		}()
	}
	time.Sleep(2 * time.Second)
}

func TestWithRedisLock(t *testing.T) {
	rd := orm.InitRedis(&orm.RedisConfig{
		Addrs:    []string{"host:6379"},
		Password: "password",
		DB:       0,
	})
	ctx := context.Background()
	g := WithRedisLock[*Single](rd)
	// 启动10个协程去请求数据，最终执行的只有1个，输出一次 I am return Single
	for i := 0; i < 10; i++ {
		go func() {
			v, err := g.Do(ctx, "key", time.Second*15, func() (*Single, error) {
				xlog.Warn("working produce Single")
				s := &Single{
					Name: "lady gaga",
					Age:  18,
				}
				time.Sleep(3 * time.Second)
				return s, nil
			})
			if err != nil {
				t.Errorf("Do error = %v", err)
			}
			xlog.Warnf("final result: %#v", v)
		}()
	}
	time.Sleep(20 * time.Second)
}
