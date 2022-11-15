package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
	"github.com/streadway/amqp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	opts        options
	channel     *amqp.Channel
	log         *logger.Logger
	conn        *amqp.Connection
	ctx         context.Context
	cancel      context.CancelFunc
	isConnected bool
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connector() *common.Connector {
	return Connector()
}

func (c *Client) getTLSConfig() (*tls.Config, error) {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: c.opts.insecure,
	}

	if c.opts.caCert != "" {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(c.opts.caCert)) {
			return nil, fmt.Errorf("error loading Root CA Cert")
		}
		tlsCfg.RootCAs = caCertPool
		c.log.Infof("TLS CA Cert Loaded for RabbitMQ Connection")
	}
	if c.opts.clientCertificate != "" && c.opts.clientKey != "" {
		cert, err := tls.X509KeyPair([]byte(c.opts.clientCertificate), []byte(c.opts.clientKey))
		if err != nil {
			return nil, fmt.Errorf("error loading tls client key pair, %s", err.Error())
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
		c.log.Infof("TLS Client Key Pair Loaded for RabbitMQ Connection")

	}
	return tlsCfg, nil
}

func (c *Client) connect() error {
	if strings.HasPrefix(c.opts.url, "amqps://") {
		tlsCfg, err := c.getTLSConfig()
		if err != nil {
			return err
		}
		c.conn, err = amqp.DialTLS(c.opts.url, tlsCfg)
		if err != nil {
			return fmt.Errorf("error dialing rabbitmq, %w", err)
		}

	} else {
		var err error
		c.conn, err = amqp.Dial(c.opts.url)
		if err != nil {
			return fmt.Errorf("error dialing rabbitmq, %w", err)
		}
	}
	var err error
	c.channel, err = c.conn.Channel()
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("error getting rabbitmq channel, %w", err)
	}
	c.isConnected = true
	return nil
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
	if err := c.connect(); err != nil {
		return err
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
	md["content_type"] = msg.ContentType
	md["content_encoding"] = msg.ContentEncoding
	md["routing_key"] = msg.RoutingKey

	headers, err := json.MarshalToString(msg.Headers)
	if err == nil {
		md["headers"] = headers
	}
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing rabbitmq metadata, %s", err.Error())
	}
	return str
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if !c.isConnected {
		if err := c.connect(); err != nil {
			return err
		}
	}
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
					channel := delivery.RoutingKey
					if delivery.Exchange != "" {
						channel = fmt.Sprintf("%s.%s", delivery.Exchange, delivery.RoutingKey)
					}
					req.SetChannel(channel)
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
