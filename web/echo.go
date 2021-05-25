package web

import (
	"fmt"
	"strings"

	"github.com/iGoogle-ink/gotil/limit"
	"github.com/labstack/echo/v4"
)

type EchoEngine struct {
	Echo *echo.Echo
	addr string
}

func InitEcho(c *Config) *EchoEngine {
	e := echo.New()

	if c.Limit != nil && c.Limit.Rate != 0 {
		e.Use(limit.NewLimiter(c.Limit).EchoLimit())
	}

	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		c.Port = ":" + c.Port
	}

	engine := &EchoEngine{Echo: e, addr: c.Host + c.Port}
	return engine
}

// Release
func (e *EchoEngine) Release() *EchoEngine {
	e.Echo.Debug = false
	e.Echo.HideBanner = true
	e.Echo.HidePort = true
	return e
}

func (e *EchoEngine) Start() {
	go func() {
		if err := e.Echo.Start(e.addr); err != nil {
			panic(fmt.Sprintf("web server port(%s) run error(%+v).", e.addr, err))
		}
	}()
}
