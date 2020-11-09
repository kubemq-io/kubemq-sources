package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"google.golang.org/api/option"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	name   string
	opts   options
	client *pubsub.Client
	log    *logger.Logger
}

func New() *Client {
	return &Client{}

}
func (c *Client) Connector() *common.Connector {
	return Connector()
}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	b := []byte(c.opts.credentials)
	client, err := pubsub.NewClient(ctx, c.opts.projectID, option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

func (c *Client) createMetadataString(msg *pubsub.Message) string {
	md := map[string]string{}
	md["id"] = msg.ID
	if len(msg.Attributes) > 0 {
		a, err := json.Marshal(msg.Attributes)
		if err != nil {
			return fmt.Sprintf("error parsing Attributes, %s", err.Error())
		}
		md["attributes"] = fmt.Sprintf("%s", a)
	}
	md["publish_time"] = msg.PublishTime.String()
	if msg.DeliveryAttempt != nil {
		md["delivery_attempt"] = fmt.Sprintf("%d", *msg.DeliveryAttempt)
	}
	md["ordering_key"] = msg.OrderingKey
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing pubsub.message metadata, %s", err.Error())
	}
	return str
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {

	var receivedMessage = make(chan *pubsub.Message, 1)

	var errCh = make(chan error, 1)

	sub := c.client.Subscription(c.opts.subscriberID)

	go func() {
		for {
			select {
			case msg := <-receivedMessage:
				b, err := json.Marshal(msg)
				if err != nil {
					c.log.Errorf("failed to parse message on error %s", err.Error())
				}
				req := types.NewRequest().SetMetadata(c.createMetadataString(msg)).SetData(b)
				_, err = target.Do(ctx, req)
				if err != nil {
					msg.Nack()
					c.log.Errorf("error processing message %s", err.Error())
				} else {
					msg.Ack()
				}
			case <-ctx.Done():
				return
			case err := <-errCh:
				c.log.Errorf("error failed to pull message on error %s", err.Error())
			}
		}
	}()

	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			receivedMessage <- msg
		})
		if err != nil {
			errCh <- err
		}
	}()

	return nil
}

func (c *Client) Stop() error {
	return c.client.Close()
}
