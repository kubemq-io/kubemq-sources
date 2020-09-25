package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	targetMiddleware "github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

const (
	gracefulShutdown = 30 * time.Second
)

var (
	errInvalidTarget = errors.New("invalid target received, cannot be nil")
)

type Server struct {
	name   string
	opts   options
	target targetMiddleware.Middleware
	log    *logger.Logger
	echo   *echo.Echo
}

func New() *Server {
	return &Server{}

}

func (s *Server) Init(ctx context.Context, cfg config.Spec) error {
	s.name = cfg.Name
	s.log = logger.NewLogger(cfg.Name)
	var err error
	s.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	s.echo = echo.New()
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())
	s.echo.POST(s.opts.path, s.process)
	s.echo.PUT(s.opts.path, s.process)

	return nil
}

func (s *Server) Start(ctx context.Context, target targetMiddleware.Middleware) error {
	if target == nil {
		return errInvalidTarget
	} else {
		s.target = target
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.echo.Start(fmt.Sprintf("%s:%d", s.opts.host, s.opts.port))
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Second):
		return nil
	}

}

func (s *Server) process(c echo.Context) error {
	req := types.NewRequest()
	err := c.Bind(req)
	if err != nil {
		return c.JSON(400, types.NewResponse().SetMetadataKeyValue("error", err.Error()))
	}
	resp, err := s.target.Do(c.Request().Context(), req)
	if err != nil {
		return c.JSON(200, types.NewResponse().SetMetadataKeyValue("error", err.Error()))
	}
	return c.JSON(200, resp)
}
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdown)
	defer cancel()
	return s.echo.Shutdown(ctx)
}
