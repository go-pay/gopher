package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/ecode"
	"github.com/go-pay/limiter"
	"github.com/go-pay/gopher/web"
	"github.com/go-pay/xlog"
	"github.com/go-pay/xtime"
)

type MemStats struct {
	Alloc        string `json:"alloc"`
	TotalAlloc   string `json:"total_alloc"`
	Sys          string `json:"sys"`
	HeapAlloc    string `json:"heap_alloc"`
	HeapSys      string `json:"heap_sys"`
	HeapIdle     string `json:"heap_idle"`
	HeapInuse    string `json:"heap_inuse"`
	HeapReleased string `json:"heap_released"`
	Frees        string `json:"frees"`
	StackInuse   string `json:"stack_inuse"`
	StackSys     string `json:"stack_sys"`
	GcSys        string `json:"gc_sys"`
	NextGc       string `json:"next_gc"`
	LastGc       string `json:"last_gc"`
	NumGc        int    `json:"num_gc"`
	NumForcedGc  int    `json:"num_forced_gc"`
	EnableGc     bool   `json:"enable_gc"`
}

func main1() {
	c := &web.Config{
		Addr:         ":2233",
		Debug:        false,
		ReadTimeout:  xtime.Duration(15 * time.Second),
		WriteTimeout: xtime.Duration(15 * time.Second),
		Limiter: &limit.Config{
			Rate:       0, // 0 速率不限流
			BucketSize: 100,
		},
	}

	g := web.InitGin(c)
	//g.Gin.Use( /*g.CORS(),*/)

	xlog.Level = xlog.DebugLevel
	ecode.Success = ecode.NewV2(0, "SUCCESS", "成功")
	initRoute(g.Gin)

	// add hook
	g.AddShutdownHook(func(c context.Context) {
		sec := 0
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			<-ticker.C
			sec++
			xlog.Warnf("second: %ds", sec)
		}
	}).AddExitHook(func(c context.Context) {
		xlog.Warn("after close hook1")
	}, func(c context.Context) {
		xlog.Warn("after close hook2")
	})
	// start server
	g.Start()
}

func initRoute(g *gin.Engine) {
	g.GET("/a/:abc", func(c *gin.Context) {
		xlog.Debug(c.Param("abc"))
		xlog.Debug(c.Request.RequestURI)
		rsp := &struct {
			Param string `json:"param"`
			Path  string `json:"path"`
		}{Param: c.Param("abc"), Path: c.Request.RequestURI}
		web.JSON(c, rsp, nil)
	})
	g.GET("/b", func(c *gin.Context) {
		web.JSON(c, nil, ecode.UnauthorizedErr)
	})
	g.POST("/c", func(c *gin.Context) {
		body, err := web.ReadRequestBody(c.Request)
		if err != nil {
			xlog.Error(err)
			web.JSON(c, nil, err)
			return
		}
		xlog.Debugf("body:%s", body)
		var ss = struct {
			Name string `json:"name"`
		}{}
		_ = c.ShouldBindJSON(&ss)
		web.JSON(c, ss, nil)
	})
	g.GET("/d", func(c *gin.Context) {
		web.JSON(c, web.Pager{PageNo: 1, PageSize: 15}.Apply(30, "我是15条数据"), nil)
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

	// proxy
	g.GET("/gopher/web/memStats", memStats)

	// postman request: GET http://localhost:2233/proxy/a
	g.GET("/proxy/a", func(c *gin.Context) {
		rsp, err := web.GinProxy[*MemStats](c, "", "http://localhost:2233", "/gopher/web/memStats")
		if err != nil {
			xlog.Errorf("GinProxy err: %v", err)
			web.JSON(c, nil, err)
			return
		}
		web.JSON(c, rsp, nil)
	})
}

func memStats(c *gin.Context) {
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	rsp := &struct {
		HeapAlloc    string
		HeapIdle     string
		HeapInuse    string
		HeapReleased string
		Frees        string
		GCSys        string
		NextGC       string
		LastGC       string
		NumGC        uint32
		NumForcedGC  uint32
		EnableGC     bool
	}{
		HeapAlloc:    fmt.Sprintf("%d(MB)", ms.HeapAlloc/1024/1024),
		HeapIdle:     fmt.Sprintf("%d(MB)", ms.HeapIdle/1024/1024),
		HeapInuse:    fmt.Sprintf("%d(MB)", ms.HeapInuse/1024/1024),
		HeapReleased: fmt.Sprintf("%d(MB)", ms.HeapReleased/1024/1024),
		Frees:        fmt.Sprintf("%d(MB)", ms.Frees/1024/1024),
		GCSys:        fmt.Sprintf("%d(MB)", ms.GCSys/1024/1024),
		NextGC:       fmt.Sprintf("%d(MB)", ms.NextGC/1024/1024),
		LastGC:       time.Unix(0, int64(ms.LastGC)).Format("2006-01-02 15:04:05.000"),
		NumGC:        ms.NumGC,
		NumForcedGC:  ms.NumForcedGC,
		EnableGC:     ms.EnableGC,
	}
	web.JSON(c, rsp, nil)
}
