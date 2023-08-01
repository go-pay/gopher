package limit

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/limit/group"
	"github.com/go-pay/gopher/limit/rate"
)

var (
	defaultConfig = &Config{
		Rate:       1000,
		BucketSize: 1000,
	}
)

type Config struct {
	// per second request，0 不限流
	Rate int `json:"rate" yaml:"rate" toml:"rate"`

	// max size，桶内最大量
	BucketSize int `json:"bucket_size" yaml:"bucket_size" toml:"bucket_size"`
}

// 速率限制器
type RateLimiter struct {
	C            *Config
	LimiterGroup *group.RateGroup
}

func NewLimiter(c *Config) (rl *RateLimiter) {
	if c == nil {
		c = defaultConfig
	}
	rl = &RateLimiter{
		C: c,
		LimiterGroup: group.NewRateGroup(func() any {
			return newLimiter(c)
		}),
	}
	return rl
}

func (r *RateLimiter) GinLimit(appName ...string) gin.HandlerFunc {
	var limitKey string
	if len(appName) > 0 {
		limitKey = appName[0]
	}
	return func(c *gin.Context) {
		if limitKey == "" {
			limitKey = strings.Split(c.Request.RequestURI, "?")[0][1:]
		}
		// log.Warning("key:", path[1:])
		limiter, ok := r.LimiterGroup.Get(limitKey).(*rate.Limiter)
		if ok && limiter != nil {
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
				return
			}
			c.Next()
		}
	}
}

func newLimiter(c *Config) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(c.Rate), c.BucketSize)
}
