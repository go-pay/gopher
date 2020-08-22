package web

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iGoogle-ink/gotil/limit"
	"github.com/iGoogle-ink/gotil/xlog"
)

func TestInitServer(t *testing.T) {
	c := &Config{
		Port: ":2233",
		Limit: &limit.Config{
			Rate:       0, // 0 速率不限流
			BucketSize: 100,
		},
	}

	g := InitServer(c)
	initRoute(g.Gin)
	g.Start()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			xlog.Warningf("get a signal %s, stop the process", si.String())
			// todo something
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func initRoute(g *gin.Engine) {
	g.GET("/a/:abc", func(c *gin.Context) {
		xlog.Debug(c.Param("abc"))
		xlog.Debug(c.Request.RequestURI)
		rsp := &struct {
			Param string `json:"param"`
			Path  string `json:"path"`
		}{Param: c.Param("abc"), Path: c.Request.RequestURI}
		JSON(c, rsp, nil)
	})
	g.GET("/b", func(c *gin.Context) {
		JSON(c, "b", nil)
	})
	g.GET("/c", func(c *gin.Context) {
		JSON(c, "c", nil)
	})
	g.GET("/d", func(c *gin.Context) {
		JSON(c, "d", nil)
	})
}
