package middleware

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
)

type LogLevelType string

const (
	LogLevelTypeNoLog LogLevelType = ""
	LogLevelTypeDebug LogLevelType = "debug"
	LogLevelTypeInfo  LogLevelType = "info"
	LogLevelTypeError LogLevelType = "error"
)

var logLevelMap = map[string]string{
	"debug": "debug",
	"info":  "info",
	"error": "error",
	"":      "",
}

type LogMiddleware struct {
	minLevel LogLevelType
	*logger.Logger
}

func NewLogMiddleware(name string, meta types.Metadata) (*LogMiddleware, error) {
	lm := &LogMiddleware{
		minLevel: LogLevelTypeNoLog,
		Logger:   logger.NewLogger(name),
	}
	level, err := meta.ParseStringMap("log_level", logLevelMap)
	if err != nil {
		return nil, fmt.Errorf("invalid log level value, %w", err)
	}
	switch level {
	case "debug":
		lm.minLevel = LogLevelTypeDebug
	case "info":
		lm.minLevel = LogLevelTypeInfo
	case "error":
		lm.minLevel = LogLevelTypeError
	}
	return lm, nil
}
