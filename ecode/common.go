package ecode

import "net/http"

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = http.StatusInternalServerError
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// ClientClosed is non-standard http status code,
	// which defined by nginx.
	// https://httpstatus.in/499/
	ClientClosed = 499
)

var (
	// base error
	Success               = New(http.StatusOK, "SUCCESS", "success")
	RequestErr            = New(http.StatusBadRequest, "PARAM_ERROR", "request param error")
	UnauthorizedErr       = New(http.StatusUnauthorized, "SIGN_ERROR", "sign error")
	ForbiddenErr          = New(http.StatusForbidden, "NO_AUTH", "no auth")
	NotFoundErr           = New(http.StatusNotFound, "RESOURCE_NOT_FOUND", "resource not found")
	TooManyRequestErr     = New(http.StatusTooManyRequests, "RATELIMIT_EXCEEDED", "ratelimit exceeded")
	ServerErr             = New(http.StatusInternalServerError, "SERVER_ERROR", "server error")
	BadGatewayErr         = New(http.StatusBadGateway, "SERVICE_UNAVAILABLE", "service offline, unavailable")
	ServiceUnavailableErr = New(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "service protected, unavailable")
)
