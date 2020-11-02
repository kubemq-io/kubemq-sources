package eventhubs

import (
	"context"
	"fmt"
	"github.com/Azure/azure-event-hubs-go/v3"
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
	client *eventhub.Hub
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
	c.log = logger.NewLogger(c.name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client, err = eventhub.NewHubFromConnectionString(c.opts.connectionString)
	if err != nil {
		return fmt.Errorf("error connecting to eventhub at %s: %w", c.opts.connectionString, err)
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
	if c.opts.partitionID != "" {
		_, err := c.client.Receive(ctx, c.opts.partitionID, func(ctx context.Context, event *eventhub.Event) error {
			go c.processIncomingMessages(ctx, event, c.opts.partitionID, errCh)
			return nil
		})
		if err != nil {
			return fmt.Errorf("error subscription to eventhubs destination on partitionID %s, %w", c.opts.partitionID, err)
		}
	}
	
	return nil
}

func (c *Client) processIncomingMessages(ctx context.Context, event *eventhub.Event, partitionID string, errCh chan error) error {
	req := types.NewRequest().SetData(event.Data)
	_, err := c.target.Do(ctx, req)
	if err != nil {
		errCh <- fmt.Errorf("error processing eventhubs eventID %s and partitionID %s , error:%s", event.ID, partitionID, err.Error())
	}
	return nil
}

func (c *Client) Stop() error {
	return c.client.Close(context.Background())
}
