package query

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
)

type Client struct {
	name   string
	opts   options
	client *kubemq.Client
	log    *logger.Logger

}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
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
	queryMetadata, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	queryResponse, err := c.client.Q().
		SetId(queryMetadata.id).
		SetTimeout(queryMetadata.timeout).
		SetChannel(queryMetadata.channel).
		SetMetadata(queryMetadata.metadata).
		SetBody(request.Data).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("error", queryResponse.Error).
			SetMetadataKeyValue("query_id", queryResponse.QueryId).
			SetMetadataKeyValue("response_client_id", queryResponse.ResponseClientId).
			SetMetadataKeyValue("executed", fmt.Sprintf("%t", queryResponse.Executed)).
			SetMetadataKeyValue("executed_at", fmt.Sprintf("%s", queryResponse.ExecutedAt)).
			SetData(queryResponse.Body),
		nil
}
