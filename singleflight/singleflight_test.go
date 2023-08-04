package singleflight

import (
	"github.com/go-pay/gopher/xlog"
	"testing"
	"time"
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
			v, err := g.Do("key", func() (*Single, error) {
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
