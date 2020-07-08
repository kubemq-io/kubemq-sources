package command

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
	cmdMetadata, err := parseMetadata(request.Metadata, c.opts)
	if err != nil {
		return nil, err
	}
	cmdResponse, err := c.client.C().
		SetId(cmdMetadata.id).
		SetTimeout(cmdMetadata.timeout).
		SetChannel(cmdMetadata.channel).
		SetMetadata(cmdMetadata.metadata).
		SetBody(request.Data).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return types.NewResponse().
			SetMetadataKeyValue("error", cmdResponse.Error).
			SetMetadataKeyValue("command_id", cmdResponse.CommandId).
			SetMetadataKeyValue("response_client_id", cmdResponse.ResponseClientId).
			SetMetadataKeyValue("executed", fmt.Sprintf("%t", cmdResponse.Executed)).
			SetMetadataKeyValue("executed_at", fmt.Sprintf("%s", cmdResponse.ExecutedAt)),
		nil
}
