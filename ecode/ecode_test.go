package ecode

import (
	"testing"

	"github.com/go-pay/gopher/xlog"
)

func TestEcode(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	e := AnalyseError(ForbiddenErr)
	xlog.Debug(e.Error())
	xlog.Debug(e.Code())
	xlog.Debug(e.Message())
	xlog.Info("============================")

	e2 := AnalyseError(UnauthorizedErr)
	xlog.Debug(e2.Error())
	xlog.Debug(e2.Code())
	xlog.Debug(e2.Message())
	xlog.Info("============================")

	sms := New(10000, "中国电信")
	xlog.Debug(sms.Error())
	xlog.Debug(sms.Code())
	xlog.Debug(sms.Message())
	xlog.Info("============================")

	mms := New(10086, "中国移动")
	xlog.Debug(mms.Error())
	xlog.Debug(mms.Code())
	xlog.Debug(mms.Message())
}
