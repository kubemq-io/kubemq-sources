package rabbitmq

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/streadway/amqp"
)

type Client struct {
	name    string
	opts    options
	channel *amqp.Channel
	log     *logger.Logger
	conn    *amqp.Connection
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
	c.conn, err = amqp.Dial(c.opts.url)
	if err != nil {
		return fmt.Errorf("error dialing rabbitmq, %w", err)
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("error getting rabbitmq channel, %w", err)
	}

	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	deliveries, err := c.channel.Consume(
		c.opts.queue,
		c.opts.consumer,
		c.opts.autoAck,
		c.opts.exclusive,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error on queue Consume: %w", err)
	}

	go func() {
		errCh := c.conn.NotifyClose(make(chan *amqp.Error))
		for {
			select {
			case delivery := <-deliveries:
				req := types.NewRequest().SetData(delivery.Body)
				_, err := target.Do(ctx, req)
				if err != nil {
					c.log.Errorf("error processing request %s", err.Error())
					_ = delivery.Reject(c.opts.requeueOnError)
				} else {
					_ = delivery.Ack(false)
				}
			case <-ctx.Done():
				return
			case err := <-errCh:
				c.log.Errorf("error on rabbitmq connection, %w", err.Reason)
			}
		}
	}()
	return nil
}

func (c *Client) Stop() error {
	return c.channel.Close()
}
