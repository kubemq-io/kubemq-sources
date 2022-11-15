package queue

import (
	"context"
	"fmt"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"

	"github.com/kubemq-io/kubemq-go/queues_stream"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *queues_stream.QueuesStreamClient
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}

func (c *Client) Init(ctx context.Context, cfg config.Spec, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}

	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client, err = queues_stream.NewQueuesStreamClient(ctx,
		queues_stream.WithAddress(c.opts.host, c.opts.port),
		queues_stream.WithClientId(c.opts.clientId),
		queues_stream.WithCheckConnection(true),
		queues_stream.WithAutoReconnect(true),
		queues_stream.WithAuthToken(c.opts.authToken),
		queues_stream.WithConnectionNotificationFunc(
			func(msg string) {
				c.log.Infof(msg)
			}),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *Client) getChannel(request *types.Request) string {
	if request.Channel != "" {
		return fmt.Sprintf("%s%s", c.opts.channel, request.Channel)
	}
	return c.opts.channel
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	queueMessage := queues_stream.NewQueueMessage().
		SetChannel(c.getChannel(request)).
		SetMetadata(request.Metadata).
		SetBody(request.Data).
		SetPolicyDelaySeconds(c.opts.delaySeconds).
		SetPolicyExpirationSeconds(c.opts.expirationSeconds).
		SetPolicyMaxReceiveCount(c.opts.maxReceiveCount).
		SetPolicyMaxReceiveQueue(c.opts.deadLetterQueue)
	result, err := c.client.Send(ctx, queueMessage)
	if err != nil {
		return nil, err
	}
	if len(result.Results) > 0 {
		if result.Results[0].IsError {
			return nil, fmt.Errorf(result.Results[0].Error)
		}
	}
	return types.NewResponse(), nil
}
