package ecode

import "sync"

var (
	errorMap = new(sync.Map)
	// base error
	OK              = New(0, "success")
	RequestErr      = New(10400, "request error")
	InvalidSignErr  = New(10401, "invalid sign")
	InvalidAppidErr = New(10402, "invalid appid")
	InvalidTokenErr = New(10403, "invalid token")
	NothingFound    = New(10404, "nothing found")
	ServerErr       = New(10500, "server error")
)
