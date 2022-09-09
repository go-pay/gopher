package web

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/ecode"
	"github.com/go-pay/gopher/xlog"
	"github.com/go-pay/gopher/xtime"
)

// CORS gin middleware cors
func (g *GinEngine) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin == "" {
			origin = c.Request.Host
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			// 允许跨域返回的Header
			c.Header("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, Session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
			// 允许的方法
			c.Header("Access-Control-Allow-Methods", "POST, PUT ,GET, OPTIONS, DELETE, HEAD, TRACE, UPDATE")
			// 允许客户端解析的Header
			c.Header("Access-Control-Expose-Headers", "Authorization, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// 缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			// 允许客户端传递校验信息，cookie
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Request.Header.Del("Origin")
		c.Next()
	}
}

// Recovery gin middleware recovery
func (g *GinEngine) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			rawReq, body []byte
		)
		body, _ = ReadRequestBody(c.Request)
		if c.Request != nil {
			rawReq, _ = httputil.DumpRequest(c.Request, true)
		}
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				stack := make([]byte, size)
				stack = stack[:runtime.Stack(stack, false)]
				bs, _ := json.Marshal(RecoverInfo{
					Time:        time.Now().Format(xtime.TimeLayout_1),
					RequestURI:  c.Request.Host + c.Request.RequestURI,
					Body:        string(body),
					RequestInfo: string(rawReq),
					Err:         err,
					Stack:       string(stack),
				})
				xlog.Errorf("[GinPanic] %s", string(bs))
				c.AbortWithError(http.StatusInternalServerError, ecode.ServerErr)
			}
		}()
		c.Next()
	}
}
