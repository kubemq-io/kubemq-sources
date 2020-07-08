package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/kubemq-hub/kubemq-source-connectors/pkg/logger"
	"github.com/kubemq-hub/kubemq-source-connectors/targets"
	"github.com/kubemq-hub/kubemq-source-connectors/types"
	"github.com/kubemq-io/kubemq-go"
)

type Client struct {
	name   string
	opts   options
	client *kubemq.Client
	log    *logger.Logger
	target targets.Target
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client, _ = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken))

	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	queueMetadata, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	queueMessage := c.client.NewQueueMessage().
		SetId(queueMetadata.id).
		SetChannel(queueMetadata.channel).
		SetMetadata(queueMetadata.metadata).
		SetBody(request.Data).
		SetPolicyDelaySeconds(queueMetadata.delaySeconds).
		SetPolicyExpirationSeconds(queueMetadata.expirationSeconds).
		SetPolicyMaxReceiveCount(queueMetadata.maxReceiveCount).
		SetPolicyMaxReceiveQueue(queueMetadata.deadLetterQueue)
	result, err := queueMessage.Send(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("error", result.Error).
			SetMetadataKeyValue("queue_id", queueMetadata.id).
			SetMetadataKeyValue("delayed_to", fmt.Sprintf("%d", result.DelayedTo)).
			SetMetadataKeyValue("expire_at", fmt.Sprintf("%d", result.ExpirationAt)).
			SetMetadataKeyValue("is_error", fmt.Sprintf("%t", result.IsError)),
		nil
}
