package web

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iGoogle-ink/gopher/xlog"
)

func TestInitServer(t *testing.T) {
	// 需要测试请自行解开注释测试

	//c := &Config{
	//	Port: ":2233",
	//	Limit: &limit.Config{
	//		Rate:       0, // 0 速率不限流
	//		BucketSize: 100,
	//	},
	//	Trace: &trace.Config{
	//		ServiceName: "trace-demo",
	//		Endpoint:    "",
	//	},
	//}
	//
	//g := InitGin(c)
	//g.Gin.Use(g.CORS(), g.Recovery())
	//
	//initRoute(g.Gin)
	//
	//g.Start()
	//
	//ch := make(chan os.Signal)
	//signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	//for {
	//	si := <-ch
	//	switch si {
	//	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
	//		xlog.Warnf("get a signal %s, stop the process", si.String())
	//		// todo something close
	//		g.Close()
	//		return
	//	case syscall.SIGHUP:
	//	default:
	//		return
	//	}
	//}
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
		time.Sleep(1 * time.Second)
		JSON(c, "b", nil)
	})
	g.POST("/c", func(c *gin.Context) {
		all, _ := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		JSON(c, string(all), nil)
	})
	g.GET("/d", func(c *gin.Context) {
		JSON(c, Pager{PageNo: 1, PageSize: 15}.Apply(30, "我是15条数据"), nil)
	})
}
