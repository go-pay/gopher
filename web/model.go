package web

import (
	"github.com/iGoogle-ink/gopher/limit"
	"github.com/iGoogle-ink/gopher/trace"
)

type Config struct {
	// http host
	Host string `json:"host" yaml:"host" toml:"host"`

	// http export port. :8080
	Port string `json:"port" yaml:"port" toml:"port"`

	// interface limit
	Limit *limit.Config `json:"limit" yaml:"limit" toml:"limit"`

	// jaeger trace config
	Trace *trace.Config `json:"trace" yaml:"trace" toml:"trace"`
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
