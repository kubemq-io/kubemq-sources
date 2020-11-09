package eventhubs

import (
	"context"
	"fmt"
	"github.com/Azure/azure-event-hubs-go/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

func (c *Client) createMetadataString(event *eventhub.Event) string {
	md := map[string]string{}
	if event.PartitionKey != nil {
		md["partition_key"] = *event.PartitionKey
	}
	if len(event.Properties) > 0 {
		a, err := json.Marshal(event.Properties)
		if err != nil {
			return fmt.Sprintf("error parsing Properties, %s", err.Error())
		}
		md["properties"] = fmt.Sprintf("%s", a)
	}
	md["id"] = event.ID
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing eventhub.Event metadata, %s", err.Error())
	}
	return str
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
		}, c.opts.receiveType)
		if err != nil {
			return fmt.Errorf("error subscription to eventhubs destination on partitionID %s, %w", c.opts.partitionID, err)
		}
	}

	return nil
}

func (c *Client) processIncomingMessages(ctx context.Context, event *eventhub.Event, partitionID string, errCh chan error) {
	req := types.NewRequest().SetMetadata(c.createMetadataString(event)).SetData(event.Data)
	_, err := c.target.Do(ctx, req)
	if err != nil {
		errCh <- fmt.Errorf("error processing eventhubs eventID %s and partitionID %s , error:%s", event.ID, partitionID, err.Error())
	}
}

func (c *Client) Stop() error {
	return c.client.Close(context.Background())
}
