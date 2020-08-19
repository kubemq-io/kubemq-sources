package api

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/binding"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

type Server struct {
	echoWebServer  *echo.Echo
	bindingService *binding.Service
}

func Start(ctx context.Context, port int, bs *binding.Service) (*Server, error) {
	s := &Server{
		echoWebServer:  echo.New(),
		bindingService: bs,
	}
	s.echoWebServer.Use(middleware.Recover())
	s.echoWebServer.Use(middleware.CORS())

	s.echoWebServer.GET("/health", func(c echo.Context) error {

		return c.String(200, "ok")

	})
	s.echoWebServer.GET("/ready", func(c echo.Context) error {

		return c.String(200, "ready")

	})
	s.echoWebServer.GET("/metrics", echo.WrapHandler(s.bindingService.PrometheusHandler()))
	s.echoWebServer.GET("/stats", func(c echo.Context) error {

		return c.JSONPretty(200, s.bindingService.Stats(), "\t")

	})
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.echoWebServer.Start(fmt.Sprintf("0.0.0.0:%d", port))
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
		return s, nil
	case <-time.After(1 * time.Second):
		return s, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("error strarting api server, %w", ctx.Err())
	}
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = s.echoWebServer.Shutdown(ctx)
}
