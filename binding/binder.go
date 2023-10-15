package binding

import (
	"context"
	"fmt"

	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/pkg/metrics"
	"github.com/kubemq-io/kubemq-sources/sources"
	"github.com/kubemq-io/kubemq-sources/sources/http"
	"github.com/kubemq-io/kubemq-sources/targets"
)

type Binder struct {
	name              string
	log               *logger.Logger
	source            sources.Source
	target            targets.Target
	md                middleware.Middleware
	httpSourceHandler *http.Client
}

func NewBinder() *Binder {
	return &Binder{}
}

func (b *Binder) buildMiddleware(cfg config.BindingConfig, exporter *metrics.Exporter) (middleware.Middleware, error) {
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
	md := middleware.Chain(b.target, middleware.RateLimiter(rateLimiter), middleware.Retry(retry), middleware.Metric(met))
	return md, nil
}

func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig, exporter *metrics.Exporter, logLevel string) error {
	b.name = cfg.Name

	b.log = logger.NewLogger(fmt.Sprintf("binding-%s", cfg.Name), logLevel)
	var err error
	b.target, err = targets.Init(ctx, cfg.Target, cfg.Name, b.log)
	if err != nil {
		return fmt.Errorf("error loading target conntector on binding %s, %w", b.name, err)
	}
	b.log.Infof("binding: %s target: initialized successfully", b.name)
	b.md, err = b.buildMiddleware(cfg, exporter)
	if err != nil {
		return fmt.Errorf("error loading middlewares on binding %s, %w", b.name, err)
	}
	b.source, err = sources.Init(ctx, cfg.Source, b.log)
	if err != nil {
		return fmt.Errorf("error loading source conntector on binding %s, %w", b.name, err)
	}
	b.log.Infof("binding: %s, source: %s, initialized successfully", b.name, cfg.Source.Name)
	if cfg.Source.Kind == "http" {
		val, ok := b.source.(*http.Client)
		if ok {
			b.httpSourceHandler = val
		}
	}
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
	err = b.target.Stop()
	if err != nil {
		return err
	}
	b.log.Infof("binding: %s, stopped successfully", b.name)
	return nil
}
