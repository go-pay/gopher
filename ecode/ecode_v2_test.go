package ecode

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-pay/gopher/xlog"
)

var (
	// base error
	V1ParameterErr = New(1000400, "request param error")
)

func TestEcodeWithReason(t *testing.T) {
	xlog.Level = xlog.DebugLevel
	e := FromError(V1ParameterErr)
	xlog.Debug(e.Error())   // error: code = 1000400 reason =  message = equest param error metadata = map[] cause = <nil>
	xlog.Debug(e.Code())    // 1000400
	xlog.Debug(e.Message()) // request param error
	xlog.Info("============================")

	Success = NewV2(1, "SUCCESS", "success")

	e2 := FromError(nil)
	xlog.Debug(e2.Error())   // error: code = 1 reason = SUCCESS message = success metadata = map[] cause = <nil>
	xlog.Debug(e2.Code())    // 1
	xlog.Debug(e2.Reason())  // SUCCESS
	xlog.Debug(e2.Message()) // success
	xlog.Info("============================")

	sms := NewV2(10000, "CTCC", "中国电信").WithMetadata(map[string]string{
		"name":   "jerry",
		"reason": "我是metadata",
	})
	xlog.Debug(sms.Error())   // error: code = 10000 reason = CTCC message = 中国电信 metadata = map[name:jerry reason:我是metadata] cause = <nil>
	xlog.Debug(sms.Code())    // 10000
	xlog.Debug(sms.Reason())  // CTCC
	xlog.Debug(sms.Message()) // 中国电信
	xlog.Debug(sms.Metadata)  // map[name:jerry reason:我是metadata]
	xlog.Info("============================")

	mms := NewV2(10086, "CMCC", "中国移动").WithCause(errors.New("我是原因"))
	xlog.Debug(mms.Error())   // error: code = 10086 reason = CMCC message = 中国移动 metadata = map[] cause = 我是原因
	xlog.Debug(mms.Code())    // 10086
	xlog.Debug(mms.Reason())  // CMCC
	xlog.Debug(mms.Message()) // 中国电信
	xlog.Debug(mms.Unwrap())  // 我是原因
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
