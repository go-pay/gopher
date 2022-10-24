package ecode

import (
	"math"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error struct
type Error struct {
	code    int    `json:"code"`
	message string `json:"message"`
}

// New error code and msg
func New(code int, msg string) *Error {
	return &Error{code: code, message: msg}
}

func (e *Error) Error() string {
	return e.message
}

// Code return error code
func (e *Error) Code() int { return e.code }

// Message return error message
func (e *Error) Message() string {
	return e.message
}

func (e *Error) GRPCStatus() *status.Status {
	return status.New(codes.Code(uint32(math.Abs(float64(e.Code())))), e.Error())
}

// analyse error info
func AnalyseError(err error) *Error {
	if err == nil {
		return New(http.StatusOK, "success")
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return errStringToError(err.Error())
}

func errStringToError(e string) *Error {
	if e == "" {
		return New(http.StatusOK, "success")
	}
	i, err := strconv.Atoi(e)
	if err != nil {
		return New(-1, e)
	}
	return New(i, e)
}
