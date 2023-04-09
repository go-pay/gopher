package smap

import (
	"testing"

	"github.com/go-pay/gopher/xlog"
)

type SmapTest struct {
	Name string
	Age  int
}

func TestSmap(t *testing.T) {
	xlog.Level = xlog.InfoLevel
	sm := Map[string, *SmapTest]{}

	actual, loaded := sm.LoadOrStore("test", &SmapTest{Name: "test", Age: 10})
	if !loaded {
		xlog.Warnf("LoadOrStore not have and store: %v", actual)
	}
	xlog.Infof("actual: %v", actual)
	value, ok := sm.Load("test")
	if !ok {
		xlog.Errorf("Load not have")
	}
	xlog.Infof("value: %v", value)

	// ==================================

	andDelete, ok := sm.LoadAndDelete("test")
	if !ok {
		xlog.Errorf("LoadAndDelete not have")
		return
	}
	xlog.Infof("andDelete: %v", andDelete)

	_, ok = sm.Load("test")
	xlog.Infof("after load and delete load sm[test] is %v", ok)

	sm.Store("test2", &SmapTest{Name: "test2", Age: 20})

	v2, ok := sm.Load("test2")
	if !ok {
		xlog.Errorf("sm[test2] Load not have")
		return
	}
	xlog.Infof("sm[test2] value: %v", v2)
}
