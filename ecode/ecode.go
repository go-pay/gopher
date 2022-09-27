package ecode

import (
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// Error struct
type Error struct {
	Status
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v cause = %v", e.Status.Code, e.Status.Reason, e.Status.Message, e.Metadata, e.cause)
}

// Code returns the code of the error.
func (e *Error) Code() int { return int(e.Status.Code) }

// Message returns the message of the error.
func (e *Error) Message() string { return e.Status.Message }

// Reason returns the reason of the error.
func (e *Error) Reason() string { return e.Status.Reason }

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Status.Code == e.Status.Code && se.Status.Reason == e.Status.Reason
	}
	return false
}

// GRPCStatus returns the Status represented by error.
func (e *Error) GRPCStatus() *status.Status {
	gs, _ := status.New(DefaultConverter.ToGRPCCode(int(e.Status.Code)), e.Status.Message).
		WithDetails(&errdetails.ErrorInfo{Reason: e.Status.Reason, Metadata: e.Metadata})
	return gs
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := DeepClone(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := DeepClone(e)
	err.Metadata = md
	return err
}

// ============================================================================================================

// New returns an error object for the code, message.
func New(code int, reason string, message string) *Error {
	return &Error{Status: Status{
		Code:    int32(code),
		Reason:  reason,
		Message: message,
	}}
}

// DeepClone deep clone error to a new error.
func DeepClone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		cause: err.cause,
		Status: Status{
			Code:     err.Status.Code,
			Reason:   err.Status.Reason,
			Message:  err.Status.Message,
			Metadata: metadata,
		},
	}
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return New(UnknownCode, err.Error(), UnknownReason)
	}
	ret := New(DefaultConverter.FromGRPCCode(gs.Code()), gs.Message(), UnknownReason)
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			ret.Status.Reason = d.Reason
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

// analyse error info
func AnalyseError(err error) *Error {
	if err == nil {
		return Success
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return errStringToError(err.Error())
}

func errStringToError(e string) *Error {
	if e == "" {
		return Success
	}
	i, err := strconv.Atoi(e)
	if err != nil {
		return New(-1, UnknownReason, e)
	}
	return New(i, UnknownReason, e)
}
