package xlog

const (
	ErrorLevel LogLevel = iota + 1
	WarnLevel
	InfoLevel
	DebugLevel
)

const (
	callDepth = 3
)

type LogLevel int

var (
	debugLog XLogger = &DebugLogger{}
	infoLog  XLogger = &InfoLogger{}
	warnLog  XLogger = &WarnLogger{}
	errLog   XLogger = &ErrorLogger{}

	Level LogLevel
)

type XLogger interface {
	LogOut(col *ColorType, format *string, args ...any)
}

func Info(args ...any) {
	infoLog.LogOut(nil, nil, args...)
}

func Infof(format string, args ...any) {
	infoLog.LogOut(nil, &format, args...)
}

func Debug(args ...any) {
	debugLog.LogOut(nil, nil, args...)
}

func Debugf(format string, args ...any) {
	debugLog.LogOut(nil, &format, args...)
}

func Warn(args ...any) {
	warnLog.LogOut(nil, nil, args...)
}

func Warnf(format string, args ...any) {
	warnLog.LogOut(nil, &format, args...)
}

func Error(args ...any) {
	errLog.LogOut(nil, nil, args...)
}

func Errorf(format string, args ...any) {
	errLog.LogOut(nil, &format, args...)
}
