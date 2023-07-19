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
	hookMaps         map[hookType][]func(c context.Context)
}

func InitGin(c *Config) *GinEngine {
	if c == nil {
		c = &Config{Addr: ":2233"}
	}
	g := gin.New()
	engine := &GinEngine{Gin: g, addrPort: c.Addr, hookMaps: make(map[hookType][]func(c context.Context))}

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

// 添加 GinServer 服务启动时的钩子函数
func (g *GinEngine) AddStartHook(hooks ...HookFunc) *GinEngine {
	for _, fn := range hooks {
		if fn != nil {
			g.hookMaps[_HookStart] = append(g.hookMaps[_HookStart], fn)
		}
	}
	return g
}

// 添加 GinServer 服务关闭时的钩子函数
func (g *GinEngine) AddCloseHook(hooks ...HookFunc) *GinEngine {
	for _, fn := range hooks {
		if fn != nil {
			g.hookMaps[_HookClose] = append(g.hookMaps[_HookClose], fn)
		}
	}
	return g
}

// 添加 GinServer 进程退出时钩子函数
func (g *GinEngine) AddExitHook(hooks ...HookFunc) *GinEngine {
	for _, fn := range hooks {
		if fn != nil {
			g.hookMaps[_HookExit] = append(g.hookMaps[_HookExit], fn)
		}
	}
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
	// call when server start hooks
	for _, fn := range g.hookMaps[_HookStart] {
		fn(context.Background())
	}
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
			// close gin http server
			g.Close()
			ctx, cancelFunc := context.WithTimeout(context.Background(), g.timeout)
			// call before close hooks
			go func() {
				if a := recover(); a != nil {
					xlog.Errorf("panic: %v", a)
				}
				for _, fn := range g.hookMaps[_HookClose] {
					fn(ctx)
				}
			}()
			// wait for a second to finish processing
			time.Sleep(g.timeout)
			cancelFunc()
			// call after close hooks
			for _, fn := range g.hookMaps[_HookExit] {
				fn(context.Background())
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
