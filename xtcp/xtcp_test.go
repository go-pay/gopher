package xtcp

import (
	"fmt"
	"github.com/go-pay/gopher/xlog"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	xlog.Level = xlog.InfoLevel
	msgCount := 0
	server, err := NewServer(&Config{
		Host:      "",
		Port:      "1122",
		HeartBeat: 30 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	server.OnConnect(func(ctx *Context) {
		xlog.Infof("OnConnect: %+v", ctx)
	})
	server.OnClose(func(ctx *Context) {
		xlog.Infof("OnClose: %+v", ctx)
	})
	server.OnMessage(func(ctx *Context) {
		xlog.Infof("OnMessage: %+v", ctx)
		if msgCount%5 == 0 {
			ctx.Conn().SendText(fmt.Sprintf("hello world, I am server %d", msgCount))
		}
		msgCount++
	})
	server.Run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			xlog.Color(xlog.Yellow).Warnf("get a signal %s, stop the process", si.String())
			// wait for program finish processing
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func TestNewClient(t *testing.T) {
	xlog.Level = xlog.InfoLevel
	client, err := NewClient(&Config{
		Host:      "",
		Port:      "1122",
		HeartBeat: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	client.OnConnect(func(ctx *Context) {
		xlog.Infof("OnConnect: %+v", ctx)
	})
	client.OnClose(func(ctx *Context) {
		xlog.Infof("OnClose: %+v", ctx)
	})
	client.OnMessage(func(ctx *Context) {
		xlog.Infof("OnMessage: %+v", ctx)
		ctx.Conn().SendText("I am Client, I get " + ctx.Message())
	})
	client.Run()

	go func() {
		for {
			xlog.Warnf("client send message")
			if err := client.SendText("hello world"); err != nil {
				xlog.Errorf("client send msg err: %v", err)
			}
			time.Sleep(time.Second * 3)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			xlog.Color(xlog.Yellow).Warnf("get a signal %s, stop the process", si.String())
			// wait for program finish processing
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
