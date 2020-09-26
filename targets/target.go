package targets

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/targets/command"
	"github.com/kubemq-hub/kubemq-sources/targets/events"
	event_store "github.com/kubemq-hub/kubemq-sources/targets/events-store"
	"github.com/kubemq-hub/kubemq-sources/targets/query"
	"github.com/kubemq-hub/kubemq-sources/types"
)

type Target interface {
	Init(ctx context.Context, cfg config.Spec) error
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
}

func Init(ctx context.Context, cfg config.Spec) (Target, error) {

	switch cfg.Kind {
	case "target.command":
		target := command.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.query":
		target := query.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events":
		target := events.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events-store":
		target := event_store.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil

	default:
		return nil, fmt.Errorf("invalid kind %s for target %s", cfg.Kind, cfg.Name)
	}

}
