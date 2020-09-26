package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
)

type Client struct {
	opts   options
	client *kubemq.Client
	sendCh chan *kubemq.Event
}

func New() *Client {
	return &Client{}

}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken),
		kubemq.WithCheckConnection(false),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	queueMessage := c.client.NewQueueMessage().
		SetChannel(c.opts.channel).
		SetMetadata(request.Metadata.String()).
		SetBody(request.Data).
		SetPolicyDelaySeconds(c.opts.delaySeconds).
		SetPolicyExpirationSeconds(c.opts.expirationSeconds).
		SetPolicyMaxReceiveCount(c.opts.maxReceiveCount).
		SetPolicyMaxReceiveQueue(c.opts.deadLetterQueue)
	result, err := queueMessage.Send(ctx)
	if err != nil {
		return nil, err
	}
	if result.IsError {
		return nil, fmt.Errorf(result.Error)
	}
	return types.NewResponse(), nil
}
