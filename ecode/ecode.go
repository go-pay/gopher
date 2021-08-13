package ecode

import (
	"math"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// eError struct
type eError struct {
	code    int    `json:"code"`
	message string `json:"message"`
}

// New new error code and msg
func New(code int, msg string) *eError {
	return &eError{code: code, message: msg}
}

func (e *eError) Error() string {
	return e.message
}

// Code return error code
func (e *eError) Code() int { return e.code }

// Message return error message
func (e *eError) Message() string {
	return e.message
}

func (e *eError) GRPCStatus() *status.Status {
	return status.New(codes.Code(uint32(math.Abs(float64(e.Code())))), e.Error())
}

// analyse error info
func AnalyseError(err error) *eError {
	if err == nil {
		return OK
	}
	if e, ok := err.(*eError); ok {
		return e
	}
	return errStringToError(err.Error())
}

func errStringToError(e string) *eError {
	if e == "" {
		return OK
	}
	i, err := strconv.Atoi(e)
	if err != nil {
		return New(-1, e)
	}
	return New(i, e)
}
