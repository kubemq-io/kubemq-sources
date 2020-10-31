package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/pkg/metrics"
	httpsrc "github.com/kubemq-hub/kubemq-sources/sources/http"
	"net/http"
	"sync"
	"time"
)

const (
	addRetryInterval = 1 * time.Second
)

type Service struct {
	sync.Mutex
	bindings          map[string]*Binder
	log               *logger.Logger
	exporter          *metrics.Exporter
	currentCtx        context.Context
	currentCancelFunc context.CancelFunc
}

func New() (*Service, error) {
	s := &Service{
		Mutex:    sync.Mutex{},
		bindings: make(map[string]*Binder),
		log:      logger.NewLogger("binding-service"),
	}
	var err error
	s.exporter, err = metrics.NewExporter()
	if err != nil {
		return nil, fmt.Errorf("failed to to initialized metrics exporter, %w", err)
	}
	return s, nil
}
func (s *Service) Start(ctx context.Context, cfg *config.Config) error {
	s.currentCtx, s.currentCancelFunc = context.WithCancel(ctx)
	if len(cfg.Bindings) == 0 {
		return nil
	}
	wg := sync.WaitGroup{}
	wg.Add(len(cfg.Bindings))
	for _, bindingCfg := range cfg.Bindings {
		go func(ctx context.Context, cfg config.BindingConfig) {
			defer wg.Done()
			err := s.Add(ctx, cfg)
			if err == nil {
				return
			} else {
				s.log.Errorf("failed to initialized binding, %s", err.Error())
			}
			count := 0
			for {
				select {
				case <-time.After(addRetryInterval):
					count++
					err := s.Add(ctx, cfg)
					if err != nil {
						s.log.Errorf("failed to initialized binding: %s, attempt: %d, error: %s", cfg.Name, count, err.Error())
					} else {
						return
					}
				case <-ctx.Done():
					return
				}
			}

		}(s.currentCtx, bindingCfg)

	}
	wg.Wait()
	return nil

}

func (s *Service) Stop() {
	for _, binder := range s.bindings {
		err := s.Remove(binder.name)
		if err != nil {
			s.log.Error(err)
		}
	}
	s.currentCancelFunc()
}
func (s *Service) Add(ctx context.Context, cfg config.BindingConfig) error {
	s.Lock()
	defer s.Unlock()
	binder := NewBinder()
	err := binder.Init(ctx, cfg, s.exporter)
	if err != nil {
		return err
	}

	err = binder.Start(ctx)
	if err != nil {
		return err
	}
	s.bindings[cfg.Name] = binder
	return nil
}

func (s *Service) Remove(name string) error {
	s.Lock()
	defer s.Unlock()
	binder, ok := s.bindings[name]
	if !ok {
		return fmt.Errorf("binding %s no found", name)
	}
	err := binder.Stop()
	if err != nil {
		return err
	}
	delete(s.bindings, name)
	return nil
}

func (s *Service) PrometheusHandler() http.Handler {
	return s.exporter.PrometheusHandler()
}
func (s *Service) Stats() []*metrics.Report {
	return s.exporter.Store.List()
}

func (s *Service) GetHttpHandlers() []*httpsrc.Handler {
	var list []*httpsrc.Handler
	s.Lock()
	defer s.Unlock()
	for _, binder := range s.bindings {
		if binder.httpSourceHandler != nil {
			list = append(list, binder.httpSourceHandler)
		}
	}
	return list
}
