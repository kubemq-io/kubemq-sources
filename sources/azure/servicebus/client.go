package servicebus

import (
	"context"
	"fmt"
	"github.com/Azure/azure-service-bus-go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
)

type Client struct {
	name   string
	opts   options
	log    *logger.Logger
	client *servicebus.Queue
	target middleware.Middleware
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(c.opts.connectionString))
	if err != nil {
		return err
	}
	c.client, err = ns.NewQueue(c.opts.queueName)
	if err != nil {
		return fmt.Errorf("error connecting to servicebus at %s: %w", c.opts.connectionString, err)
	}
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	} else {
		c.target = target
	}
	errCh := make(chan error, 1)
	go func() {
		for {
			select {
			case err := <-errCh:
				c.log.Errorf(err.Error())
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		err := c.client.Receive(ctx, servicebus.HandlerFunc(func(ctx context.Context, message *servicebus.Message) error {
			go c.processIncomingMessages(ctx, message, errCh)
			return nil
		}))
		if err != nil {
			return
		}
	}()

	return nil
}

func (c *Client) processIncomingMessages(ctx context.Context, message *servicebus.Message, errCh chan error) {
	req := types.NewRequest().SetData(message.Data)
	_, err := c.target.Do(ctx, req)
	if err != nil {
		errCh <- fmt.Errorf("error processing servicebus messageID %s, error:%s", message.ID, err.Error())
	}
}

func (c *Client) Stop() error {
	return c.client.Close(context.Background())
}
