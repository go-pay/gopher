package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
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
	wg               sync.WaitGroup
	addrPort         string
	IgnoreReleaseLog bool
	hookMaps         map[hookType][]func(c context.Context)
}

func InitGin(c *Config) *GinEngine {
	if c == nil {
		c = &Config{Addr: ":2233"}
	}
	g := gin.New()
	engine := &GinEngine{Gin: g, wg: sync.WaitGroup{}, addrPort: c.Addr, hookMaps: make(map[hookType][]func(c context.Context))}

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

// 添加 GinServer 服务关闭时的钩子函数
func (g *GinEngine) AddShutdownHook(hooks ...HookFunc) *GinEngine {
	for _, fn := range hooks {
		if fn != nil {
			g.hookMaps[_HookShutdown] = append(g.hookMaps[_HookShutdown], fn)
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
	// monitoring signal
	go g.goNotifySignal()

	// start gin http server
	xlog.Warnf("Listening and serving HTTP on %s", g.addrPort)
	if err := g.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			xlog.Warn("http: Server closed")
		} else {
			panic(fmt.Sprintf("server.ListenAndServe(), error(%+v).", err))
		}
	}
	xlog.Color(xlog.Green).Warnf("wait for process working finished")
	// wait for process finished
	g.wg.Wait()
	xlog.Color(xlog.Green).Warnf("process exit")
}

// 监听信号
func (g *GinEngine) goNotifySignal() {
	g.wg.Add(1)
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
				for _, fn := range g.hookMaps[_HookShutdown] {
					fn(ctx)
				}
			}()
			// wait for program finish processing
			xlog.Warnf("waiting for the process to finish %v", g.timeout)
			time.Sleep(g.timeout)
			cancelFunc()
			// call after close hooks
			for _, fn := range g.hookMaps[_HookExit] {
				fn(context.Background())
			}
			// notify process exit
			g.wg.Done()
			runtime.Gosched()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
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
