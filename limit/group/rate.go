package group

import (
	"sync"
)

// RateGroup 速率Group
type RateGroup struct {
	new func() any
	rgs sync.Map
	sync.RWMutex
}

// NewRateGroup 新建RateGroup
func NewRateGroup(new func() any) (rg *RateGroup) {
	if new == nil {
		panic("RateGroup: can't assign a nil to the new function")
	}
	return &RateGroup{new: new}
}

// Get 获取RateGroup，如果没有则新建
func (r *RateGroup) Get(key string) any {
	rg, ok := r.rgs.Load(key)
	if !ok {
		r.RLock()
		newRg := r.new
		r.RUnlock()
		rg = newRg()
		r.rgs.Store(key, rg)
	}
	return rg
}
