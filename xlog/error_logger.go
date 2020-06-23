package xlog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
)

type ErrorLogger struct {
	logger *log.Logger
	once   sync.Once
}

func (i *ErrorLogger) logOut(format *string, v ...interface{}) {
	i.once.Do(func() {
		i.init()
	})
	if format != nil {
		i.logger.Output(3, fmt.Sprintf(*format, v...))
		//i.logger.Writer().Write(stack())
		return
	}
	i.logger.Output(3, fmt.Sprintln(v...))
	//i.logger.Writer().Write(stack())
}

func (i *ErrorLogger) init() {
	i.logger = log.New(os.Stderr, "[ERROR] >> ", log.Lmsgprefix|log.Lshortfile|log.Lmicroseconds|log.Ldate)
}

func stack() (bs []byte) {
	var buf [2 << 10]byte
	runtime.Stack(buf[:], false)
	return buf[:]
}
