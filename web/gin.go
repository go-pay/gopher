package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iGoogle-ink/gopher/limit"
	"github.com/iGoogle-ink/gopher/trace"
	"github.com/iGoogle-ink/gopher/xlog"
	"github.com/iGoogle-ink/gopher/xtime"
)

type GinEngine struct {
	server *http.Server
	Gin    *gin.Engine
	Tracer *trace.Tracer
	addr   string
}

func InitGin(c *Config) *GinEngine {
	g := gin.New()
	engine := &GinEngine{Gin: g, addr: c.Addr}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = xtime.Duration(60 * time.Second)
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = xtime.Duration(60 * time.Second)
	}

	engine.server = &http.Server{
		Addr:         engine.addr,
		Handler:      g,
		ReadTimeout:  time.Duration(c.ReadTimeout),
		WriteTimeout: time.Duration(c.WriteTimeout),
	}
	g.Use(gin.Logger(), engine.Recovery())
	if c.Trace != nil {
		engine.Tracer = trace.NewTracer(c.Trace)
		g.Use(engine.Tracer.GinTrace())
	}
	if c.Limit != nil && c.Limit.Rate != 0 {
		g.Use(limit.NewLimiter(c.Limit).GinLimit())
	}
	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	return engine
}

func (g *GinEngine) Start() {
	go func() {
		xlog.Debugf("Listening and serving HTTP on %s", g.addr)
		if err := g.server.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("web server addr(%s) run error(%+v).", g.addr, err))
		}
	}()
}

func (g *GinEngine) Close() {
	if g.Tracer != nil {
		g.Tracer.Close()
	}
	if g.server != nil {
		_ = g.server.Shutdown(context.Background())
	}
}
