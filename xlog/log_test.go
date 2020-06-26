package xlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	Debug("debug")
	Error("error")
	Info("info")
	Warning("warning")
}
