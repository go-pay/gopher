package web

import (
	"context"

	"github.com/go-pay/gopher/limit"
	"github.com/go-pay/gopher/trace"
	"github.com/go-pay/gopher/xtime"
)

const (
	_HookShutdown hookType = "shutdown"
	_HookExit     hookType = "exit"
)

type hookType string

type HookFunc func(c context.Context)

type Config struct {
	Addr         string         `json:"addr" yaml:"addr" toml:"addr"`                            // addr, default :2233
	ReadTimeout  xtime.Duration `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`    // read_timeout, default 60s
	WriteTimeout xtime.Duration `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"` // write_timeout, default 60s
	Debug        bool           `json:"debug" yaml:"debug" toml:"debug"`                         // is show log
	Limiter      *limit.Config  `json:"limiter" yaml:"limiter" toml:"limiter"`                   // interface limit
	Trace        *trace.Config  `json:"trace" yaml:"trace" toml:"trace"`                         // jaeger trace config
}

type RecoverInfo struct {
	Time        string `json:"time"`
	RequestURI  string `json:"request_uri"`
	Body        string `json:"body"`
	RequestInfo string `json:"request_info"`
	Err         any    `json:"error"`
	Stack       string `json:"stack"`
}

type CommonRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type HttpRsp[V any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    V      `json:"data,omitempty"`
}
