package ws

import "errors"

const (
	CommonChanCount = 1000
)

var (
	ErrConnClosed = errors.New("connection is closed")
)
