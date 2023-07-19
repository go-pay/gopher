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
	server              *http.Server
	Gin                 *gin.Engine
	Tracer              *trace.Tracer
	timeout             time.Duration
	addrPort            string
	IgnoreReleaseLog    bool
	beforeCloseHookFunc []func()
	afterCloseHookFunc  []func()
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

// 注册 GinServer 关闭前的钩子函数
func (g *GinEngine) RegBeforeCloseHooks(hooks ...func()) *GinEngine {
	g.beforeCloseHookFunc = append(g.beforeCloseHookFunc, hooks...)
	return g
}

// 注册 GinServer Pod 关闭后的钩子函数
func (g *GinEngine) RegAfterCloseHooks(hooks ...func()) *GinEngine {
	g.afterCloseHookFunc = append(g.afterCloseHookFunc, hooks...)
	return g
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
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			xlog.Color(xlog.Yellow).Warnf("get a signal %s, stop the process", si.String())
			// call before close hooks
			for _, fn := range g.beforeCloseHookFunc {
				fn()
			}
			// close gin http server
			g.Close()
			// wait for a second to finish processing
			time.Sleep(g.timeout)
			// call after close hooks
			for _, fn := range g.afterCloseHookFunc {
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
