package sources

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/kubemq-hub/kubemq-source-connectors/sources/kubemq/query"
	"github.com/kubemq-hub/kubemq-source-connectors/sources/kubemq/queue"
	"github.com/kubemq-hub/kubemq-source-connectors/targets"
)

type Source interface {
	Init(ctx context.Context, cfg config.Metadata) error
	Start(ctx context.Context, target targets.Target) error
	Stop() error
	Name() string
}

func Load(ctx context.Context, cfgs []config.Metadata) (map[string]Source, error) {
	sources := make(map[string]Source)
	for _, cfg := range cfgs {
		_, ok := sources[cfg.Name]
		if ok {
			return nil, fmt.Errorf("duplicated source name found, %s", cfg.Name)
		}

		switch cfg.Kind {
		case "source.queue":
			source := queue.New()
			if err := source.Init(ctx, cfg); err != nil {
				return nil, err
			}
			sources[cfg.Name] = source
		case "source.rpc":
			source := query.New()
			if err := source.Init(ctx, cfg); err != nil {
				return nil, err
			}
			sources[cfg.Name] = source

		default:
			return nil, fmt.Errorf("invalid source kind %s for source %s", cfg.Kind, cfg.Name)
		}
	}
	return sources, nil
}
