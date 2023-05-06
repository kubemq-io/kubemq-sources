package events

import (
	"context"
	"fmt"
	"time"

	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
)

const (
	defaultSendTimeout     = 10 * time.Second
	defaultStreamReconnect = 1 * time.Second
)

type Client struct {
	log    *logger.Logger
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

func (c *Client) Init(ctx context.Context, cfg config.Spec, bindingName string, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger(cfg.Kind)
	}
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	c.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(fmt.Sprintf("kubemq-sources_%s_%s", bindingName, c.opts.clientId)),
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
		return fmt.Sprintf("%s%s", c.opts.channel, request.Channel)
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
