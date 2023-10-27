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
	Success               = NewV2(http.StatusOK, "SUCCESS", "success")
	RequestErr            = NewV2(http.StatusBadRequest, "PARAM_ERROR", "request param error")
	UnauthorizedErr       = NewV2(http.StatusUnauthorized, "SIGN_ERROR", "sign error")
	ForbiddenErr          = NewV2(http.StatusForbidden, "NO_AUTH", "no auth")
	NotFoundErr           = NewV2(http.StatusNotFound, "RESOURCE_NOT_FOUND", "resource not found")
	TooManyRequestErr     = NewV2(http.StatusTooManyRequests, "RATELIMIT_EXCEEDED", "ratelimit exceeded")
	ServerErr             = NewV2(http.StatusInternalServerError, "SERVER_ERROR", "server error")
	BadGatewayErr         = NewV2(http.StatusBadGateway, "SERVICE_UNAVAILABLE", "service offline, unavailable")
	ServiceUnavailableErr = NewV2(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "service protected, unavailable")
)

var (
	SuccessV1               = New(http.StatusOK, "success")
	RequestErrV1            = New(http.StatusBadRequest, "request param error")
	UnauthorizedErrV1       = New(http.StatusUnauthorized, "sign error")
	ForbiddenErrV1          = New(http.StatusForbidden, "no auth")
	NotFoundErrV1           = New(http.StatusNotFound, "resource not found")
	TooManyRequestErrV1     = New(http.StatusTooManyRequests, "ratelimit exceeded")
	ServerErrV1             = New(http.StatusInternalServerError, "server error")
	BadGatewayErrV1         = New(http.StatusBadGateway, "service offline, unavailable")
	ServiceUnavailableErrV1 = New(http.StatusServiceUnavailable, "service protected, unavailable")
)
