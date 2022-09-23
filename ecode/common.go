package ecode

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// ClientClosed is non-standard http status code,
	// which defined by nginx.
	// https://httpstatus.in/499/
	ClientClosed = 499
)

var (
	// base error
	Success         = New(0, "success")
	RequestErr      = New(10400, "request error")
	InvalidSignErr  = New(10401, "invalid sign")
	InvalidAppidErr = New(10402, "invalid appid")
	InvalidTokenErr = New(10403, "invalid token")
	NothingFound    = New(10404, "nothing found")
	ServerErr       = New(10500, "server error")
)
