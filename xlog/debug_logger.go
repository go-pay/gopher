package xlog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type DebugLogger struct {
	logger *log.Logger
	once   sync.Once
}

func (d *DebugLogger) LogOut(col *ColorType, format *string, v ...any) {
	d.once.Do(func() {
		d.init()
	})
	if Level >= DebugLevel {
		if col != nil {
			if format != nil {
				_ = d.logger.Output(callDepth, string(*col)+fmt.Sprintf(*format, v...)+string(Reset))
				return
			}
			_ = d.logger.Output(callDepth, string(*col)+fmt.Sprintln(v...)+string(Reset))
			return
		}
		if format != nil {
			_ = d.logger.Output(callDepth, fmt.Sprintf(*format, v...))
			return
		}
		_ = d.logger.Output(callDepth, fmt.Sprintln(v...))
	}
}

func (d *DebugLogger) init() {
	if Level == 0 {
		Level = WarnLevel
	}
	d.logger = log.New(os.Stdout, "[DEBUG] >> ", log.Lmsgprefix|log.Lshortfile|log.Ldate|log.Lmicroseconds)
}
