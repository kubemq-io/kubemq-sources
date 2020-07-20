package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/sources"
	"github.com/kubemq-hub/kubemq-sources/targets"
)

type Binder struct {
	source sources.Source
	target targets.Target
	md     middleware.Middleware
}

func New() *Binder {
	return &Binder{}
}

func (b *Binder) InitTarget(ctx context.Context, cfg config.Metadata, target targets.Target) error {
	err := target.Init(ctx, cfg)
	if err != nil {
		return err
	}
	b.target = target
	b.md = middleware.Chain(target)
	return nil
}

func (b *Binder) InitSource(ctx context.Context, cfg config.Metadata, source sources.Source) error {
	err := source.Init(ctx, cfg)
	if err != nil {
		return err
	}
	b.source = source
	return nil
}

func (b *Binder) Start(ctx context.Context) error {
	if b.md == nil {
		return fmt.Errorf("no valid initialzed target middleware found")
	}
	if b.source == nil {
		return fmt.Errorf("no valid initialzed source found")
	}

}
