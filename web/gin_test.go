package web

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/ecode"
	"github.com/go-pay/gopher/xlog"
)

func TestInitServer(t *testing.T) {
	// 需要测试请自行解开注释测试

	//c := &Config{
	//	Addr:  ":2233",
	//	Debug: true,
	//	Limiter: &limit.Config{
	//		Rate:       0, // 0 速率不限流
	//		BucketSize: 100,
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
	//		time.Sleep(time.Second)
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
		JSON(c, nil, ecode.UnauthorizedErr)
	})
	g.POST("/c", func(c *gin.Context) {
		body, err := ReadRequestBody(c.Request)
		if err != nil {
			xlog.Error(err)
			JSON(c, nil, err)
			return
		}
		xlog.Debugf("body:%s", body)
		var ss = struct {
			Name string `json:"name"`
		}{}
		_ = c.ShouldBindJSON(&ss)
		JSON(c, ss, nil)
	})
	g.GET("/d", func(c *gin.Context) {
		JSON(c, Pager{PageNo: 1, PageSize: 15}.Apply(30, "我是15条数据"), nil)
	})
	g.POST("/wechatCallback", func(c *gin.Context) {
		//notify, err := wechat.V3ParseNotify(c.Request)
		//if err != nil {
		//	xlog.Errorf("wechat.V3ParseNotify(),err:%+v", err)
		//	return
		//}
		//xlog.Debug("Id:", notify.Id)
		//xlog.Debug("EventType:", notify.EventType)
		//xlog.Debug("ResourceType:", notify.ResourceType)
		//xlog.Debug("Resource:", notify.Resource)
		//xlog.Debug("CreateTime:", notify.CreateTime)
		//xlog.Debug("Summary:", notify.Summary)
	})
}
