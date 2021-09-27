package web

import (
	"github.com/iGoogle-ink/gopher/limit"
	"github.com/iGoogle-ink/gopher/trace"
	"github.com/iGoogle-ink/gopher/xtime"
)

type Config struct {
	Addr         string         `json:"addr" yaml:"addr" toml:"addr"`                            // addr :8080
	ReadTimeout  xtime.Duration `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`    // read_timeout
	WriteTimeout xtime.Duration `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"` // write_timeout
	Debug        bool           `json:"debug" yaml:"debug" toml:"debug"`                         // is show log
	Limit        *limit.Config  `json:"limit" yaml:"limit" toml:"limit"`                         // interface limit
	Trace        *trace.Config  `json:"trace" yaml:"trace" toml:"trace"`                         // jaeger trace config
}

type RecoverInfo struct {
	Time  string      `json:"time"`
	Url   string      `json:"url"`
	Err   string      `json:"error"`
	Query interface{} `json:"query"`
	Stack string      `json:"stack"`
}

type CommonRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
