package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/targets"
	"github.com/kubemq-hub/kubemq-sources/types"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

type Client struct {
	name   string
	opts   options
	client *pubsub.Client
	log    *logger.Logger
	target -sources.Target
	doneCh chan struct{}
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
	c.doneCh = make(chan struct{})

	client, err := pubsub.NewClient(ctx, c.opts.projectID)
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

func (c *Client) Start(ctx context.Context, target -sources.Target) error {
	if target == nil {
		return errInvalidTarget
	} else {
		c.target = target
	}
	for i := 0; i < c.opts.concurrency; i++ {
		go c.run(ctx)
	}
	return nil
}

func (c *Client) run(ctx context.Context) {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		err := c.getQueueMessages(cctx)
		if err != nil {
			c.log.Error(err.Error())
			return
		}
		select {
		case <-cctx.Done():
			return
		case <-c.doneCh:
			return
		default:

		}
	}

}
func (c *Client) getQueueMessages(ctx context.Context) error {
	client, err := pubsub.NewClient(ctx, c.opts.projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClientfailed on err: %v", err)
	}

	sub := client.Subscription(c.opts.subID)

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)
	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			err := c.processQueueMessage(ctx, msg)
			if err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()

	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil {
		return fmt.Errorf("received error: %v", err)
	}
	return nil
}

func (c *Client) processQueueMessage(ctx context.Context, msg *pubsub.Message) error {
	req, err := types.ParseRequestFromGooglePubSubMessage(msg)
	if err != nil {
		return fmt.Errorf("invalid request format, %w", err)
	}
	_, err = c.target.Do(ctx, req)
	if err != nil {
		return err
	}
	return nil

}

func (c *Client) Stop() error {
	close(c.doneCh)
	return nil
}
