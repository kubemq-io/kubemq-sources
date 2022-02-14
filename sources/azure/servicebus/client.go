package servicebus

import (
	"context"
	"fmt"
	"github.com/Azure/azure-service-bus-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
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
	req := types.NewRequest().SetMetadata(c.createMetadataString(message)).SetData(message.Data)
	_, err := c.target.Do(ctx, req)
	if err != nil {
		errCh <- fmt.Errorf("error processing servicebus messageID %s, error:%s", message.ID, err.Error())
	}
}

func (c *Client) Stop() error {
	return c.client.Close(context.Background())
}

func (c *Client) createMetadataString(message *servicebus.Message) string {
	md := map[string]string{}
	md["content_type"] = message.ContentType
	md["correlation_id"] = message.CorrelationID
	md["delivery_count"] = fmt.Sprintf("%d", message.DeliveryCount)
	if message.SessionID != nil {
		md["session_id"] = *message.SessionID
	}
	md["group_sequence"] = fmt.Sprintf("%d", message.GroupSequence)
	md["id"] = message.ID
	md["label"] = message.Label
	md["reply_to"] = message.ReplyTo
	md["to"] = message.To
	md["time_to_live"] = message.TTL.String()
	if message.LockToken != nil {
		md["lock_token"] = message.LockToken.String()
	}
	if len(message.UserProperties) > 0 {
		a, err := json.Marshal(message.UserProperties)
		if err != nil {
			return fmt.Sprintf("error parsing UserProperties, %s", err.Error())
		}
		md["user_properties"] = string(a)
	}
	md["format"] = fmt.Sprintf("%d", message.Format)
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing servicebus.Message metadata, %s", err.Error())
	}
	return str
}
