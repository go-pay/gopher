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

func (l *ZapLogger) Info(args ...any) {
	if l.Sugar != nil && Level >= InfoLevel {
		l.Sugar.Info(args...)
		return
	}
	infoLog.LogOut(nil, nil, args...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	if l.Sugar != nil && Level >= InfoLevel {
		l.Sugar.Infof(format, args...)
		return
	}
	infoLog.LogOut(nil, &format, args...)
}

func (l *ZapLogger) Debug(args ...any) {
	if l.Sugar != nil && Level >= DebugLevel {
		l.Sugar.Debug(args...)
		return
	}
	debugLog.LogOut(nil, nil, args...)
}

func (l *ZapLogger) Debugf(format string, args ...any) {
	if l.Sugar != nil && Level >= DebugLevel {
		l.Sugar.Debugf(format, args...)
		return
	}
	debugLog.LogOut(nil, &format, args...)
}

func (l *ZapLogger) Warn(args ...any) {
	if l.Sugar != nil && Level >= WarnLevel {
		l.Sugar.Warn(args...)
		return
	}
	warnLog.LogOut(nil, nil, args...)
}

func (l *ZapLogger) Warnf(format string, args ...any) {
	if l.Sugar != nil && Level >= WarnLevel {
		l.Sugar.Warnf(format, args...)
		return
	}
	warnLog.LogOut(nil, &format, args...)
}

func (l *ZapLogger) Error(args ...any) {
	if l.Sugar != nil && Level >= ErrorLevel {
		l.Sugar.Error(args...)
		return
	}
	errLog.LogOut(nil, nil, args...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	if l.Sugar != nil && Level >= ErrorLevel {
		l.Sugar.Errorf(format, args...)
		return
	}
	errLog.LogOut(nil, &format, args...)
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
