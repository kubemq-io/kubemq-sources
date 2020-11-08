package events

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
	"time"
)

const (
	defaultSendTimeout     = 10 * time.Second
	defaultStreamReconnect = 1 * time.Second
)

type Client struct {
	opts   options
	client *kubemq.Client
	sendCh chan *kubemq.Event
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Stop() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
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
		// making sure that this stays false in order the http source will work correctly
		kubemq.WithCheckConnection(false),
	)
	if err != nil {
		return err
	}
	c.sendCh = make(chan *kubemq.Event, 1)
	go c.runStreamProcessing(ctx)

	return nil
}

func (c *Client) runStreamProcessing(ctx context.Context) {
	for {
		errCh := make(chan error, 1)
		go func() {
			c.client.StreamEvents(ctx, c.sendCh, errCh)
		}()
		select {
		case <-errCh:
			time.Sleep(defaultStreamReconnect)
		case <-ctx.Done():
			goto done
		}
	}
done:
}
func (c *Client) getChannel(request *types.Request) string {
	if request.Channel != "" {
		return request.Channel
	}
	return c.opts.channel
}
func (c *Client) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	event := c.client.NewEvent().
		SetChannel(c.getChannel(request)).
		SetMetadata(request.Metadata).
		SetBody(request.Data)

	select {
	case c.sendCh <- event:
	case <-time.After(defaultSendTimeout):
		return types.NewResponse().SetError(fmt.Errorf("error timeout on sending event")), nil
	case <-ctx.Done():
		return types.NewResponse().SetError(ctx.Err()), nil
	}
	return types.NewResponse(), nil
}
