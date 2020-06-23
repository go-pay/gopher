package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	Gin  *gin.Engine
	port string
}

func InitServer(port string) *Engine {
	g := gin.Default()
	g.Use(cors())
	g.Use(gin.Recovery())
	engine := &Engine{
		Gin:  g,
		port: port,
	}
	return engine
}

func (e *Engine) Start() {
	go func() {
		if err := e.Gin.Run(e.port); err != nil {
			panic(fmt.Sprintf("web server port(%s) run error(%+v).", e.port, err))
		}
	}()
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin == "" {
			origin = c.Request.Host
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-CSRF-Token, authorization, sign, appid, ts")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
