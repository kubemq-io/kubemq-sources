package targets

import (
	"context"
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/kubemq-hub/kubemq-source-connectors/types"
)

type Target interface {
	Init(ctx context.Context, cfg config.Metadata) error
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
	Name() string
}

//
//func Load(ctx context.Context, cfgs []config.Metadata) (map[string]Target, error) {
//	targets := make(map[string]Target)
//	for _, cfg := range cfgs {
//		_, ok := targets[cfg.Name]
//		if ok {
//			return nil, fmt.Errorf("duplicated target name found, %s", cfg.Name)
//		}
//
//		switch cfg.Kind {
//		case "target.http":
//			target := http.New()
//			if err := target.Init(ctx, cfg); err != nil {
//				return nil, err
//			}
//			targets[cfg.Name] = target
//		default:
//			return nil, fmt.Errorf("invalid target kind %s for target %s", cfg.Kind, cfg.Name)
//		}
//	}
//	return targets, nil
//}
