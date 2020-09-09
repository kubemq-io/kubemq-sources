package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"google.golang.org/api/option"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

type Client struct {
	name   string
	opts   options
	client *pubsub.Client
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
	b := []byte(c.opts.credentials)
	client, err := pubsub.NewClient(ctx, c.opts.projectID,option.WithCredentialsJSON(b))
	if err != nil {
		return err
	}
	c.client = client
	return nil
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
				req := types.NewRequest().SetData(b)
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
				c.log.Errorf("error failed to pull message from queue on error %s", err.Error())
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
