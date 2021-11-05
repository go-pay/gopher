package micro

import (
	"context"

	"github.com/go-pay/gopher/xlog"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"
)

func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		xlog.Infof("server[%s], method[%s], params[%+v]", req.Service(), req.Method(), req.Body())
		return fn(ctx, req, rsp)
	}
}

func LogClientWrap(c client.Client) client.Client {
	return &logClientWrapper{c}
}

type logClientWrapper struct {
	client.Client
}

func (l *logClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	xlog.Infof("client[%s], method[%s], params[%+v]", req.Service(), req.Method(), req.Body())
	return l.Client.Call(ctx, req, rsp)
}
