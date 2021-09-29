package logger

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-sources/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var ServiceLog = NewServiceLogger()
var core = initCore()
var encoderConfig = zapcore.EncoderConfig{
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
}

func initCore() zapcore.Core {

	var w zapcore.WriteSyncer
	std, _, _ := zap.Open("stderr")
	if global.EnableLogFile {
		path, _ := os.Executable()

		err := os.MkdirAll("./logs", 0660)
		if err != nil {
			panic(err.Error())
		}
		logR := &LogRotator{
			Ctx:        context.Background(),
			Filename:   filepath.Join(filepath.Dir(path), "/logs/kubemq-sources.log"),
			MaxSize:    100, // megabytes
			MaxBackups: 5,
			MaxAge:     28, //days
		}
		w = zap.CombineWriteSyncers(std, logR, ServiceLog)
	} else {
		w = zap.CombineWriteSyncers(std, ServiceLog)
	}
	enc := zapcore.NewJSONEncoder(encoderConfig)
	if global.LoggerType == "console" {
		enc = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewCore(
		enc,
		zapcore.AddSync(w),
		zapcore.DebugLevel)
}

func LogLevelToZapLevel(value string) zapcore.Level {
	switch value {
	case "debug":
		return zap.DebugLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Infof(format, v)
}

func NewLogger(name string, level ...string) *Logger {
	lvlStr := ""
	if len(level) > 0 {
		lvlStr = level[0]
	}
	zapLogger := zap.New(core, zap.IncreaseLevel(LogLevelToZapLevel(lvlStr)))
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
