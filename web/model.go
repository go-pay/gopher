package web

import (
	"github.com/iGoogle-ink/gopher/limit"
	"github.com/iGoogle-ink/gopher/trace"
)

type Config struct {
	// http host
	Host string

	// http export port. :8080
	Port string

	// interface limit
	Limit *limit.Config

	// jaeger trace config
	Trace *trace.Config
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
