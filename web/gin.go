package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopher/limit"
	"github.com/go-pay/gopher/trace"
	"github.com/go-pay/gopher/xlog"
	"github.com/go-pay/gopher/xtime"
)

type GinEngine struct {
	server           *http.Server
	Gin              *gin.Engine
	Tracer           *trace.Tracer
	timeout          time.Duration
	addrPort         string
	IgnoreReleaseLog bool
	closeHookFunc    []func()
}

func InitGin(c *Config) *GinEngine {
	if c == nil {
		c = &Config{Addr: ":2233"}
	}
	g := gin.New()
	engine := &GinEngine{Gin: g, addrPort: c.Addr}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = xtime.Duration(60 * time.Second)
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = xtime.Duration(60 * time.Second)
	}
	engine.timeout = time.Duration(c.ReadTimeout)
	engine.server = &http.Server{
		Addr:         engine.addrPort,
		Handler:      g,
		ReadTimeout:  time.Duration(c.ReadTimeout),
		WriteTimeout: time.Duration(c.WriteTimeout),
	}
	g.Use(engine.Logger(false), engine.Recovery())
	if c.Trace != nil {
		engine.Tracer = trace.NewTracer(c.Trace)
		g.Use(engine.Tracer.GinTrace())
	}
	if c.Limiter != nil && c.Limiter.Rate != 0 {
		g.Use(limit.NewLimiter(c.Limiter).GinLimit())
	}
	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	return engine
}

// 注册GinServer关闭后的钩子函数
func (g *GinEngine) RegCloseHooks(hooks ...func()) {
	for _, hook := range hooks {
		g.closeHookFunc = append(g.closeHookFunc, hook)
	}
}

func (g *GinEngine) Start() {
	go func() {
		xlog.Warnf("Listening and serving HTTP on %s", g.addrPort)
		if err := g.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				xlog.Warn("http: Server closed")
				return
			}
			panic(fmt.Sprintf("server.ListenAndServe(), error(%+v).", err))
		}
	}()
}

// 监听信号
func (g *GinEngine) NotifySignal() {
	sig := []os.Signal{syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT}
	ch := make(chan os.Signal)
	signal.Notify(ch, sig...)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			xlog.Color(xlog.Yellow).Warnf("get a signal %s, stop the process", si.String())
			g.Close()
			time.Sleep(g.timeout)
			// Close相关服务
			for _, fn := range g.closeHookFunc {
				fn()
			}
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func (g *GinEngine) StartAndNotify() {
	g.Start()
	g.NotifySignal()
}

func (g *GinEngine) Close() {
	if g.Tracer != nil {
		g.Tracer.Close()
	}
	if g.server != nil {
		// disable keep-alives on existing connections
		g.server.SetKeepAlivesEnabled(false)
		_ = g.server.Shutdown(context.Background())
	}
}
