package query

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
	"time"
)

type Client struct {
	opts   options
	client *kubemq.Client
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
		kubemq.WithAuthToken(c.opts.authToken))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {

	queryResponse, err := c.client.Q().
		SetTimeout(time.Duration(c.opts.timeoutSeconds) * time.Second).
		SetChannel(c.opts.channel).
		SetMetadata(request.Metadata.String()).
		SetBody(request.Data).
		Send(ctx)
	if err != nil {
		return types.NewResponse().SetError(err), nil
	}
	if !queryResponse.Executed {
		return types.NewResponse().SetError(fmt.Errorf(queryResponse.Error)), nil
	}
	mt, _ := types.UnmarshallMetadata(queryResponse.Metadata)
	return types.NewResponse().
			SetMetadata(mt).
			SetData(queryResponse.Body),
		nil
}
