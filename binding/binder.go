package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/pkg/metrics"
	"github.com/kubemq-hub/kubemq-sources/sources"
	"github.com/kubemq-hub/kubemq-sources/targets"
)

type Binder struct {
	name   string
	log    *logger.Logger
	source sources.Source
	target targets.Target
	md     middleware.Middleware
}

func NewBinder() *Binder {
	return &Binder{}
}
func (b *Binder) buildMiddleware(cfg config.BindingConfig, exporter *metrics.Exporter) (middleware.Middleware, error) {
	log, err := middleware.NewLogMiddleware(cfg.Name, cfg.Properties)
	if err != nil {
		return nil, err
	}
	retry, err := middleware.NewRetryMiddleware(cfg.Properties, b.log)
	if err != nil {
		return nil, err
	}
	rateLimiter, err := middleware.NewRateLimitMiddleware(cfg.Properties)
	if err != nil {
		return nil, err
	}
	met, err := middleware.NewMetricsMiddleware(cfg, exporter)
	if err != nil {
		return nil, err
	}
	md := middleware.Chain(b.target, middleware.RateLimiter(rateLimiter), middleware.Retry(retry), middleware.Metric(met), middleware.Log(log))
	return md, nil
}
func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig, exporter *metrics.Exporter) error {
	var err error
	b.name = cfg.Name
	b.log = logger.NewLogger(b.name)
	b.target, err = targets.Init(ctx, cfg.Target)
	if err != nil {
		return fmt.Errorf("error loading target conntector %s on binding %s, %w", cfg.Target.Name, b.name, err)
	}
	b.log.Infof("binding: %s, target: %s, initialized successfully", b.name, cfg.Target.Name)
	b.md, err = b.buildMiddleware(cfg, exporter)
	if err != nil {
		return fmt.Errorf("error loading middlewares %s on binding %s, %w", cfg.Target.Name, b.name, err)
	}
	b.source, err = sources.Init(ctx, cfg.Source)
	if err != nil {
		return fmt.Errorf("error loading source conntector %s on binding %s, %w", cfg.Source.Name, b.name, err)
	}
	b.log.Infof("binding: %s, source: %s, initialized successfully", b.name, cfg.Source.Name)
	b.log.Infof("binding: %s, initialized successfully", b.name)
	return nil
}

func (b *Binder) Start(ctx context.Context) error {
	if b.md == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed target middleware found", b.name)
	}
	if b.source == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed source found", b.name)
	}
	err := b.source.Start(ctx, b.md)
	if err != nil {
		return err
	}
	b.log.Infof("binding: %s, started successfully", b.name)
	return nil
}
func (b *Binder) Stop() error {
	err := b.source.Stop()
	if err != nil {
		return err
	}
	b.log.Infof("binding: %s, stopped successfully", b.name)
	return nil
}
