package rabbitmq

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/streadway/amqp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	name    string
	opts    options
	channel *amqp.Channel
	log     *logger.Logger
	conn    *amqp.Connection
	ctx     context.Context
	cancel  context.CancelFunc
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
	c.conn, err = amqp.Dial(c.opts.url)
	if err != nil {
		return fmt.Errorf("error dialing rabbitmq, %w", err)
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("error getting rabbitmq channel, %w", err)
	}
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}
func (c *Client) createMetadataString(msg amqp.Delivery) string {
	md := map[string]string{}
	md["delivery_mode"] = fmt.Sprintf("%d", msg.DeliveryMode)
	md["priority"] = fmt.Sprintf("%d", msg.Priority)
	md["correlation_id"] = msg.CorrelationId
	md["reply_to"] = msg.ReplyTo
	md["expiration"] = msg.Expiration
	md["message_id"] = msg.MessageId
	md["timestamp"] = msg.Timestamp.String()
	md["type"] = msg.Type
	md["user_id"] = msg.UserId
	md["app_id"] = msg.AppId
	md["consumer_tag"] = msg.ConsumerTag
	md["message_count"] = fmt.Sprintf("%d", msg.MessageCount)
	md["delivery_tag"] = fmt.Sprintf("%d", msg.DeliveryTag)
	md["redelivered"] = fmt.Sprintf("%t", msg.Redelivered)
	md["exchange"] = msg.Exchange
	md["routing_key"] = msg.RoutingKey

	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing mqtt metadata, %s", err.Error())
	}
	return str
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
				req := types.NewRequest().SetData(delivery.Body).SetMetadata(c.createMetadataString(delivery))
				if c.opts.dynamicMapping {
					req.SetChannel(c.opts.queue)
				}
				_, err := target.Do(ctx, req)
				if c.opts.autoAck {
					if err != nil {
						c.log.Errorf("error processing request %s", err.Error())
					}
				} else {
					if err != nil {
						c.log.Errorf("error processing request %s", err.Error())
						_ = delivery.Reject(c.opts.requeueOnError)
					} else {
						_ = delivery.Ack(false)
					}
				}
			case <-c.ctx.Done():
				return
			case <-ctx.Done():
				return
			case err := <-errCh:
				if err != nil {
					c.log.Errorf("error on rabbitmq connection, %s", err.Reason)
				}
			}
		}
	}()
	return nil
}

func (c *Client) Stop() error {
	c.cancel()
	return c.conn.Close()
}
