package ecode

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-pay/gopher/xlog"
)

func TestEcodeWithReason(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	e := AnalyseError(ForbiddenErr)
	xlog.Debug(e.Error())
	xlog.Debug(e.Code())
	xlog.Debug(e.Message())
	xlog.Info("============================")

	Success = NewV2(1, "SUCCESS", "success")

	e2 := FromError(nil)
	xlog.Debug(e2.Error())
	xlog.Debug(e2.Code())
	xlog.Debug(e2.Reason())
	xlog.Debug(e2.Message())
	xlog.Info("============================")

	sms := NewV2(10000, "CTCC", "中国电信")
	xlog.Debug(sms.Error())
	xlog.Debug(sms.Code())
	xlog.Debug(sms.Reason())
	xlog.Debug(sms.Message())
	xlog.Info("============================")

	mms := NewV2(10086, "CMCC", "中国电信")
	xlog.Debug(mms.Error())
	xlog.Debug(mms.Code())
	xlog.Debug(mms.Reason())
	xlog.Debug(mms.Message())
}

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrorV2
		err  error
		want bool
	}{
		{
			name: "true",
			e:    NewV2(404, "test", ""),
			err:  NewV2(http.StatusNotFound, "test", ""),
			want: true,
		},
		{
			name: "false",
			e:    NewV2(0, "test", ""),
			err:  errors.New("test"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.e.Is(tt.err); ok != tt.want {
				t.Errorf("ErrorV2.ErrorV2() = %v, want %v", ok, tt.want)
			}
		})
	}
}
