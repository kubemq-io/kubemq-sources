package api

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/binding"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
	"time"
)

type Server struct {
	echoWebServer  *echo.Echo
	bindingService *binding.Service
	paths          map[string]string
}

func Start(ctx context.Context, port int, bs *binding.Service) (*Server, error) {
	s := &Server{
		echoWebServer:  echo.New(),
		bindingService: bs,
		paths: map[string]string{
			"GET/health":  "GET/health",
			"GET/ready":   "GET/ready",
			"GET/metrics": "GET/metrics",
			"GET/stats":   "GET/stats",
		},
	}
	s.echoWebServer.HideBanner = true
	s.echoWebServer.Use(middleware.Recover())
	s.echoWebServer.Use(middleware.CORS())

	s.echoWebServer.GET("/health", func(c echo.Context) error {

		return c.String(200, "ok")

	})
	s.echoWebServer.GET("/ready", func(c echo.Context) error {
		return c.String(200, "ready")
	})
	s.echoWebServer.GET("/metrics", echo.WrapHandler(s.bindingService.PrometheusHandler()))
	s.echoWebServer.GET("/bindings", func(c echo.Context) error {
		return c.JSONPretty(200, s.bindingService.GetStatus(), "\t")
	})
	s.echoWebServer.GET("/bindings/stats", func(c echo.Context) error {
		return c.JSONPretty(200, s.bindingService.Stats(), "\t")
	})
	if err := s.loadHttpHandlers(); err != nil {
		return nil, err
	}
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

func (s *Server) loadHttpHandlers() error {
	for _, handler := range s.bindingService.GetHttpHandlers() {
		for _, method := range handler.Methods {
			m := strings.ToUpper(method)
			path := handler.Path
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			pathKey := fmt.Sprintf("%s%s", m, path)
			if _, ok := s.paths[pathKey]; ok {
				return fmt.Errorf("duplicate method/path founded: %s%s", m, path)
			}
			switch m {
			case "GET":
				s.echoWebServer.GET(path, echo.WrapHandler(handler))
			case "POST":
				s.echoWebServer.POST(path, echo.WrapHandler(handler))
			case "PUT":
				s.echoWebServer.PUT(path, echo.WrapHandler(handler))
			case "PATCH":
				s.echoWebServer.PATCH(path, echo.WrapHandler(handler))
			case "DELETE":
				s.echoWebServer.DELETE(path, echo.WrapHandler(handler))
			}
			s.paths[pathKey] = pathKey
		}
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.echoWebServer.Shutdown(ctx)
}
