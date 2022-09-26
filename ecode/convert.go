package ecode

import (
	"google.golang.org/grpc/codes"
)

type statusConverter struct{}

// DefaultConverter default converter.
var DefaultConverter = statusConverter{}

// ToGRPCCode converts an HTTP error code into the corresponding gRPC response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func (c statusConverter) ToGRPCCode(code int) codes.Code {
	if code <= 0 {
		return codes.Unknown
	}
	return codes.Code(code)
	//switch code {
	//case http.StatusOK:
	//	return codes.OK
	//case http.StatusBadRequest:
	//	return codes.InvalidArgument
	//case http.StatusUnauthorized:
	//	return codes.Unauthenticated
	//case http.StatusForbidden:
	//	return codes.PermissionDenied
	//case http.StatusNotFound:
	//	return codes.NotFound
	//case http.StatusConflict:
	//	return codes.Aborted
	//case http.StatusTooManyRequests:
	//	return codes.ResourceExhausted
	//case http.StatusInternalServerError:
	//	return codes.Internal
	//case http.StatusNotImplemented:
	//	return codes.Unimplemented
	//case http.StatusServiceUnavailable:
	//	return codes.Unavailable
	//case http.StatusGatewayTimeout:
	//	return codes.DeadlineExceeded
	//case ClientClosed:
	//	return codes.Canceled
	//}
	//return codes.Unknown
}

// FromGRPCCode converts a gRPC error code into the corresponding HTTP response status.
// See: https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto
func (c statusConverter) FromGRPCCode(code codes.Code) int {
	if code == codes.Unknown {
		return -1
	}
	return int(code)
	//switch code {
	//case codes.OK:
	//	return http.StatusOK
	//case codes.Canceled:
	//	return ClientClosed
	//case codes.Unknown:
	//	return http.StatusInternalServerError
	//case codes.InvalidArgument:
	//	return http.StatusBadRequest
	//case codes.DeadlineExceeded:
	//	return http.StatusGatewayTimeout
	//case codes.NotFound:
	//	return http.StatusNotFound
	//case codes.AlreadyExists:
	//	return http.StatusConflict
	//case codes.PermissionDenied:
	//	return http.StatusForbidden
	//case codes.Unauthenticated:
	//	return http.StatusUnauthorized
	//case codes.ResourceExhausted:
	//	return http.StatusTooManyRequests
	//case codes.FailedPrecondition:
	//	return http.StatusBadRequest
	//case codes.Aborted:
	//	return http.StatusConflict
	//case codes.OutOfRange:
	//	return http.StatusBadRequest
	//case codes.Unimplemented:
	//	return http.StatusNotImplemented
	//case codes.Internal:
	//	return http.StatusInternalServerError
	//case codes.Unavailable:
	//	return http.StatusServiceUnavailable
	//case codes.DataLoss:
	//	return http.StatusInternalServerError
	//}
	//return http.StatusInternalServerError
}
