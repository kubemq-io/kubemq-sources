package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"

	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	log           *logger.Logger
	opts          options
	config        *kafka.Config
	consumerGroup string
	cg            kafka.ConsumerGroup
	target        middleware.Middleware
	consumer      consumer
	ctx           context.Context
	cancel        context.CancelFunc
}

type consumer struct {
	ready          chan bool
	callback       func(ctx context.Context, request *types.Request) (*types.Response, error)
	once           sync.Once
	dynamicMapping bool
}

func (consumer *consumer) ConsumeClaim(session kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if consumer.callback != nil {
			req := &types.Request{
				Data: message.Value,
			}
			if consumer.dynamicMapping {
				req.SetChannel(message.Topic)
			}
			req.Metadata = consumer.createMetadataString(message)
			_, err := consumer.callback(session.Context(), req)
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
	c.consumerGroup = c.opts.consumerGroup

	kc := kafka.NewConfig()
	kc.Version = kafka.V2_0_0_0
	isSSL, isSASL := c.opts.parseSecurityProtocol()

	if isSASL {
		kc.Net.SASL.Enable = isSASL
		kc.Net.SASL.User = c.opts.saslUsername
		kc.Net.SASL.Password = c.opts.saslPassword
		kc.Net.SASL.Mechanism = c.opts.parseASLMechanism()
	}
	if isSSL {
		kc.Net.TLS.Enable = true
		tlsCfg := &tls.Config{
			InsecureSkipVerify: c.opts.insecure,
		}
		if c.opts.cacert != "" {
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM([]byte(c.opts.cacert)) {
				return fmt.Errorf("error loading Root CA Cert")
			}
			tlsCfg.RootCAs = caCertPool
			c.log.Infof("TLS CA Cert Loaded for Kafka Connection")
		}
		if c.opts.clientCert != "" && c.opts.clientKey != "" {
			cert, err := tls.X509KeyPair([]byte(c.opts.clientCert), []byte(c.opts.clientKey))
			if err != nil {
				return fmt.Errorf("error loading tls client key pair, %s", err.Error())
			}
			tlsCfg.Certificates = []tls.Certificate{cert}
			c.log.Infof("TLS Client Key Pair Loaded for Kafka Connection")
		}
		kc.Net.TLS.Config = tlsCfg
	}
	c.config = kc
	c.cg, err = kafka.NewConsumerGroup(c.opts.brokers, c.consumerGroup, kc)
	if err != nil {
		return err
	}
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (consumer *consumer) createMetadataString(message *kafka.ConsumerMessage) string {
	md := map[string]string{}
	md["timestamp"] = message.Timestamp.String()
	md["block_timestamp"] = message.BlockTimestamp.String()
	md["topic"] = message.Topic
	md["partition"] = fmt.Sprintf("%d", message.Partition)
	md["offset"] = fmt.Sprintf("%d", message.Offset)
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing kafka.ConsumerMessage metadata, %s", err.Error())
	}
	return str
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	}

	c.target = target

	ready := make(chan bool)
	c.consumer = consumer{
		ready:          ready,
		callback:       target.Do,
		dynamicMapping: c.opts.dynamicMapping,
	}
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ctx.Done():
				return
			default:
				err := c.cg.Consume(ctx, c.opts.topics, &(c.consumer))
				if err != nil {
					c.log.Errorf("error processing request %s", err.Error())
				}
				if ctx.Err() != nil {
					return
				}
			}
		}
	}()
	<-ready
	return nil
}

// see https://github.com/Shopify/sarama/issues/1321
func (c *Client) Stop() error {
	c.cancel()
	if c.cg != nil {
		c.config.MetricRegistry.UnregisterAll()
		return c.cg.Close()
	}
	return nil
}
