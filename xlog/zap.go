package xlog

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
	c      zap.Config
	once   sync.Once
}

var z = &ZapLogger{}

func Zap() *ZapLogger {
	z.once.Do(func() {
		z.initZap()
	})
	return z
}

func (l *ZapLogger) Info(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Info(args...)
		return
	}
	infoLog.logOut(nil, nil, args...)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Infof(format, args...)
		return
	}
	infoLog.logOut(nil, &format, args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Debug(args...)
		return
	}
	debugLog.logOut(nil, nil, args...)
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Debugf(format, args...)
		return
	}
	debugLog.logOut(nil, &format, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Warn(args...)
		return
	}
	warnLog.logOut(nil, nil, args...)
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Warnf(format, args...)
		return
	}
	warnLog.logOut(nil, &format, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Error(args...)
		return
	}
	errLog.logOut(nil, nil, args...)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	if l.Sugar != nil {
		l.Sugar.Errorf(format, args...)
		return
	}
	errLog.logOut(nil, &format, args...)
}

func (l *ZapLogger) initZap() {
	var err error
	l.c = zap.NewProductionConfig()
	l.c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l.c.EncoderConfig.EncodeTime = timeEncoder
	l.Logger, err = l.c.Build(zap.AddCallerSkip(1))
	if err != nil {
		Errorf("l.initZap(),err:%+v", err)
		return
	}
	l.Sugar = l.Logger.Sugar()
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, "2006-01-02 15:04:05.000", enc)
}

func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}
