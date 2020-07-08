package logger

import (
	"fmt"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var zapConfig = zap.Config{
	Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
	Development:       false,
	DisableCaller:     true,
	DisableStacktrace: true,
	Sampling:          nil,
	Encoding:          "json",
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	},
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
	InitialFields:    nil,
}

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Infof(format, v)
}

func NewLogger(name string) *Logger {
	zapLogger, _ := zapConfig.Build()
	l := &Logger{
		SugaredLogger: zapLogger.Sugar().With("source", name),
	}
	return l
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	l.Info(str)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	l.Debug(str)
}
func (l *Logger) Fatalf(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	l.DPanicf(str)
}

func (l *Logger) NewWith(p1, p2 string) *Logger {
	return &Logger{SugaredLogger: l.SugaredLogger.With(p1, p2)}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.Info(string(p))
	return len(p), nil
}
