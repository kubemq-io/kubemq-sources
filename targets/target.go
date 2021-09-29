package targets

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/targets/command"
	"github.com/kubemq-io/kubemq-sources/targets/events"
	event_store "github.com/kubemq-io/kubemq-sources/targets/events-store"
	"github.com/kubemq-io/kubemq-sources/targets/query"
	"github.com/kubemq-io/kubemq-sources/targets/queue"
	"github.com/kubemq-io/kubemq-sources/types"
)

type Target interface {
	Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
	Connector() *common.Connector
	Stop() error
}

func Init(ctx context.Context, cfg config.Spec, log *logger.Logger) (Target, error) {

	switch cfg.Kind {
	case "kubemq.command":
		target := command.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "kubemq.query":
		target := query.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "kubemq.events":
		target := events.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "kubemq.events-store":
		target := event_store.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "kubemq.queue":
		target := queue.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil

	default:
		return nil, fmt.Errorf("invalid kind %s for target %s", cfg.Kind, cfg.Name)
	}

}

func Connectors() common.Connectors {
	return []*common.Connector{
		queue.Connector(),
		query.Connector(),
		events.Connector(),
		event_store.Connector(),
		command.Connector(),
	}
}
