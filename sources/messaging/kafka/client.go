package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
)

type Client struct {
	name          string
	log           *logger.Logger
	opts          options
	consumerGroup string
	cg            kafka.ConsumerGroup
	target        middleware.Middleware
	consumer      consumer
}

type consumer struct {
	ready    chan bool
	callback func(ctx context.Context, request *types.Request) (*types.Response, error)
	once     sync.Once
}

func (consumer *consumer) ConsumeClaim(session kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if consumer.callback != nil {
			_, err := consumer.callback(session.Context(), &types.Request{
				Data: message.Value,
			})
			if err != nil {
				return err
			}

			session.MarkMessage(message, "")
		}
	}
	return nil
}

func (consumer *consumer) Cleanup(kafka.ConsumerGroupSession) error {
	return nil
}

func (consumer *consumer) Setup(kafka.ConsumerGroupSession) error {
	consumer.once.Do(func() {
		close(consumer.ready)
	})

	return nil
}

func New() *Client {
	return &Client{}
}

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.consumerGroup = c.opts.consumerGroup

	kc := kafka.NewConfig()
	kc.Version = kafka.V2_0_0_0

	if c.opts.saslUsername != "" {
		kc.Net.SASL.Enable = true
		kc.Net.SASL.User = c.opts.saslUsername
		kc.Net.SASL.Password = c.opts.saslPassword

		kc.Net.TLS.Enable = true
		kc.Net.TLS.Config = &tls.Config{
			ClientAuth: 0,
		}
	}

	c.cg, err = kafka.NewConsumerGroup(c.opts.brokers, c.consumerGroup, kc)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {

	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	}

	c.target = target

	ready := make(chan bool)
	c.consumer = consumer{
		ready:    ready,
		callback: target.Do,
	}

	go func() {
		defer func() {
			c.log.Debugf("Closing ConsumerGroup for topics: %v", c.opts.topics)
			err := c.cg.Close()
			if err != nil {
				c.log.Errorf("Error closing consumer group: %v", err)
			}
		}()
		c.log.Debugf("Subscribed and listening to topics: %s", c.opts.topics)

		for {
			err := c.cg.Consume(ctx, c.opts.topics, &(c.consumer))
			if err != nil {
				c.log.Errorf("error processing request %s", err.Error())
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	<-ready
	return nil
}
func (c *Client) Stop() error {

	return nil
}
