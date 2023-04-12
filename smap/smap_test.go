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
	sm.Store("test2", &SmapTest{Name: "test2", Age: 11})
	l := sm.Len()
	xlog.Infof("len = 2 ?: %d", l)
	sm.Store("test2", &SmapTest{Name: "test2", Age: 12})
	l = sm.Len()
	xlog.Infof("len = 2 ?: %d", l)

	// ==================================

	andDelete, ok := sm.LoadAndDelete("test")
	if !ok {
		xlog.Errorf("LoadAndDelete not have")
		return
	}
	xlog.Infof("andDelete: %v", andDelete)
	l = sm.Len()
	xlog.Infof("len = 1 ?: %d", l)
	sm.Store("test3", &SmapTest{Name: "test3", Age: 13})
	l = sm.Len()
	xlog.Infof("len = 2 ?: %d", l)
	_, ok = sm.Load("test")
	xlog.Infof("after load and delete load sm[test] is %v", ok)

	sm.Store("test2", &SmapTest{Name: "test2", Age: 20})

	v2, ok := sm.Load("test2")
	if !ok {
		xlog.Errorf("sm[test2] Load not have")
		return
	}
	xlog.Infof("sm[test2] value: %v", v2)
	l = sm.Len()
	xlog.Infof("len = 2 ?: %d", l)
	sm.Delete("test2")
	l = sm.Len()
	xlog.Infof("len = 1 ?: %d", l)
	sm.Delete("test3")
	l = sm.Len()
	xlog.Infof("len = 0 ?: %d", l)

}
