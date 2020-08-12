package limit

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	xrate "github.com/iGoogle-ink/gotil/rate"
	"golang.org/x/time/rate"
)

var (
	defaultConfig = &Config{
		Rate:       1000,
		BucketSize: 1000,
	}
)

type Config struct {
	// per second request，0 不限流
	Rate int

	// max size，桶内最大量
	BucketSize int
}

// 速率限制器
type RateLimiter struct {
	C            *Config
	LimiterGroup *xrate.RateGroup
}

func NewLimiter(c *Config) (rl *RateLimiter) {
	if c == nil {
		c = defaultConfig
	}
	rl = &RateLimiter{
		C: c,
		LimiterGroup: xrate.NewRateGroup(func() interface{} {
			return newLimiter(c)
		}),
	}
	return rl
}

func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.Split(c.Request.RequestURI, "?")[0]
		// log.Warning("key:", path[1:])
		limiter := r.LimiterGroup.Get(path[1:]).(*rate.Limiter)
		if allow := limiter.Allow(); !allow {
			rsp := struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{
				Code:    10503,
				Message: "服务器忙，请稍后重试...",
			}
			c.JSON(http.StatusOK, rsp)
			c.Abort()
		}
		c.Next()
	}
}

func newLimiter(c *Config) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(c.Rate), c.BucketSize)
}
