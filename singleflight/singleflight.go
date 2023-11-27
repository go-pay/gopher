// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package singleflight provides a duplicate function call suppression
// mechanism.
package singleflight // import "golang.org/x/sync/singleflight"

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-pay/gopher/redislock"
	"github.com/go-pay/gopher/util"
	"github.com/redis/go-redis/v9"
)

// errGoexit indicates the runtime.Goexit was called in
// the user given function.
var errGoexit = errors.New("runtime.Goexit was called")

// A panicError is an arbitrary value recovered from a panic
// with the stack trace during the execution of given function.
type panicError struct {
	value any
	stack []byte
}

// Error implements error interface.
func (p *panicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", p.value, p.stack)
}

func newPanicError(v any) error {
	stack := debug.Stack()

	// The first line of the stack trace is of the form "goroutine N [status]:"
	// but by the time the panic reaches Do the goroutine may no longer exist
	// and its status will have changed. Trim out the misleading line.
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &panicError{value: v, stack: stack}
}

// call is an in-flight or completed singleflight.Do call
type call[V any] struct {
	wg sync.WaitGroup

	// These fields are written once before the WaitGroup is done
	// and are only read after the WaitGroup is done.
	val V
	err error

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	dups  int
	chans []chan<- Result[V]
}

// Group represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
type Group[V any] struct {
	//rdsClient *redislock.Client   // redis lock client
	mu sync.Mutex          // protects m
	m  map[string]*call[V] // lazily initialized
}

// Result holds the results of Do, so they can be passed
// on a channel.
type Result[V any] struct {
	Val    V
	Shared bool
	Err    error
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
// The return value shared indicates whether v was given to multiple callers.
func (g *Group[V]) Do(key string, fn func() (V, error)) (v V, shared bool, err error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[V])
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, true, c.err
	}

	c := new(call[V])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.dups > 0, c.err
}

// DoChan is like Do but returns a channel that will receive the
// results when they are ready.
//
// The returned channel will not be closed.
func (g *Group[V]) DoChan(key string, fn func() (V, error)) <-chan Result[V] {
	ch := make(chan Result[V], 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[V])
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call[V]{chans: []chan<- Result[V]{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)
	return ch
}

// doCall handles the single call for a key.
func (g *Group[V]) doCall(c *call[V], key string, fn func() (v V, err error)) {
	normalReturn := false
	recovered := false

	// use double-defer to distinguish panic from runtime.Goexit,
	// more details see https://golang.org/cl/134395
	defer func() {
		// the given function invoked runtime.Goexit
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			// In order to prevent the waiting channels from being blocked forever,
			// needs to ensure that this panic cannot be recovered.
			if len(c.chans) == 0 {
				panic(e)
			}
			go panic(e)
			select {} // Keep this goroutine around so that it will appear in the crash dump.
		} else if c.err != errGoexit {
			// Normal return
			for _, ch := range c.chans {
				ch <- Result[V]{c.val, c.dups > 0, c.err}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				// Ideally, we would wait to take a stack trace until we've determined
				// whether this is a panic or a runtime.Goexit.
				//
				// Unfortunately, the only way we can distinguish the two is to see
				// whether the recover stopped the goroutine from terminating, and by
				// the time we know that, the part of the stack trace relevant to the
				// panic has been discarded.
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}

// Forget tells the singleflight to forget about a key.  Future calls
// to Do for this key will call the function rather than waiting for
// an earlier call to complete.
func (g *Group[V]) Forget(key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}

// --------------------------------------------------------------------

type RedisGroup[V any] struct {
	rdsClient *redislock.Client // redis lock client
	//mu        sync.Mutex        // protects m
	//m         map[string]*call[V] // lazily initialized
}

func WithRedisLock[V any](rd redis.Scripter) (g RedisGroup[V]) {
	return RedisGroup[V]{rdsClient: redislock.New(rd)}
}

func (g *RedisGroup[V]) Do(ctx context.Context, key string, ttl time.Duration, fn func() (V, error)) (v V, err error) {
	obtain, err := g.rdsClient.Obtain(ctx, key, ttl)
	if err != nil && !errors.Is(err, redislock.ErrNotObtained) {
		return
	}
	dataKey := key + "_data"
	// not obtain locker, is doing
	if obtain == nil {
		times := 10
		interval := ttl / time.Duration(times)
		//xlog.Warn("没抢到锁, 等待中... %v", interval)
		// 最多循环10次获取数据
		for times > 0 {
			value, err := g.rdsClient.GetData(ctx, dataKey)
			if err != nil {
				if err == redis.Nil {
					time.Sleep(interval)
					times--
					continue
				}
				return v, err
			}
			d := new(redislock.Data)
			if err = util.UnmarshalString(util.MarshalString(value), d); err != nil {
				return v, err
			}
			if err = util.UnmarshalString(d.Data, &v); err != nil {
				return v, err
			}
			return v, nil
		}
		return v, errors.New("redis lock timeout")
	}

	// obtain locker, do function
	v, err = fn()
	err = g.rdsClient.SetData(ctx, dataKey, v, err, ttl)
	_ = obtain.Release(ctx)
	return
}
