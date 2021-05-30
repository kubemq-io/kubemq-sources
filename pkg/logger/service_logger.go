package logger

import "github.com/kardianos/service"

type ServiceLogger struct {
	logger service.Logger
}

func NewServiceLogger() *ServiceLogger {
	s := &ServiceLogger{}
	return s
}

func (s *ServiceLogger) SetLogger(logger service.Logger) {
	s.logger = logger
}

func (s *ServiceLogger) Write(p []byte) (n int, err error) {
	if s.logger == nil {
		return len(p), nil
	}
	err = s.logger.Info(string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *ServiceLogger) Sync() error {
	return nil
}
