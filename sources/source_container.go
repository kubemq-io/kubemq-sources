//go:build container
// +build container

package sources

import (
	"context"
	"fmt"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/sources/aws/amazonmq"
	"github.com/kubemq-io/kubemq-sources/sources/aws/msk"
	"github.com/kubemq-io/kubemq-sources/sources/aws/s3"
	"github.com/kubemq-io/kubemq-sources/sources/aws/sqs"
	"github.com/kubemq-io/kubemq-sources/sources/azure/eventhubs"
	"github.com/kubemq-io/kubemq-sources/sources/azure/servicebus"
	"github.com/kubemq-io/kubemq-sources/sources/gcp/pubsub"
	"github.com/kubemq-io/kubemq-sources/sources/http"
	"github.com/kubemq-io/kubemq-sources/sources/messaging/activemq"
	"github.com/kubemq-io/kubemq-sources/sources/messaging/ibmmq"
	"github.com/kubemq-io/kubemq-sources/sources/messaging/kafka"
	"github.com/kubemq-io/kubemq-sources/sources/messaging/mqtt"
	"github.com/kubemq-io/kubemq-sources/sources/messaging/nats"
	"github.com/kubemq-io/kubemq-sources/sources/messaging/rabbitmq"
	"github.com/kubemq-io/kubemq-sources/sources/storage/filesystem"
	"github.com/kubemq-io/kubemq-sources/sources/storage/minio"
)

type Source interface {
	Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error
	Start(ctx context.Context, target middleware.Middleware) error
	Stop() error
	Connector() *common.Connector
}

func Init(ctx context.Context, cfg config.Spec, log *logger.Logger) (Source, error) {
	switch cfg.Kind {
	case "messaging.activemq":
		source := activemq.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "messaging.rabbitmq":
		source := rabbitmq.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "messaging.mqtt":
		source := mqtt.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "messaging.kafka":
		source := kafka.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "aws.sqs":
		source := sqs.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "aws.amazonmq":
		source := amazonmq.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "aws.msk":
		source := msk.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "aws.s3":
		source := s3.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "gcp.pubsub":
		source := pubsub.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "azure.eventhubs":
		source := eventhubs.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "azure.servicebus":
		source := servicebus.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "messaging.ibmmq":
		target := ibmmq.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "messaging.nats":
		target := nats.New()
		if err := target.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return target, nil
	case "http":
		source := http.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "storage.filesystem":
		source := filesystem.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	case "storage.minio":
		source := minio.New()
		if err := source.Init(ctx, cfg, log); err != nil {
			return nil, err
		}
		return source, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for source %s", cfg.Kind, cfg.Name)
	}
}

func Connectors() common.Connectors {
	return []*common.Connector{
		// General
		http.Connector(),
		filesystem.Connector(),
		minio.Connector(),
		rabbitmq.Connector(),
		mqtt.Connector(),
		kafka.Connector(),
		activemq.Connector(),
		ibmmq.Connector(),
		nats.Connector(),
		// AWS
		sqs.Connector(),
		amazonmq.Connector(),
		msk.Connector(),
		s3.Connector(),
		// GCP
		pubsub.Connector(),

		// Azure
		eventhubs.Connector(),
		servicebus.Connector(),
	}
}
