package mqtt

import (
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	defaultConnectTimeout = 5 * time.Second
)

type Client struct {
	opts   options
	client mqtt.Client
	log    *logger.Logger
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
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", c.opts.host))
	opts.SetUsername(c.opts.username)
	opts.SetPassword(c.opts.password)
	opts.SetClientID(c.opts.clientId)
	opts.SetConnectTimeout(defaultConnectTimeout)
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to mqtt broker, %w", token.Error())
	}
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	} else {
		c.target = target
	}

	c.client.Subscribe(c.opts.topic, byte(c.opts.qos), func(client mqtt.Client, message mqtt.Message) {
		go c.processIncomingMessages(ctx, message)
	})

	return nil
}

func (c *Client) createMetadataString(msg mqtt.Message) string {
	md := map[string]string{}
	md["message_id"] = fmt.Sprintf("%d", msg.MessageID())
	md["topic"] = msg.Topic()
	md["duplicate"] = fmt.Sprintf("%t", msg.Duplicate())
	md["qos"] = fmt.Sprintf("%d", msg.Qos())
	md["retained"] = fmt.Sprintf("%t", msg.Retained())
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing mqtt metadata, %s", err.Error())
	}
	return str
}
func (c *Client) processIncomingMessages(ctx context.Context, msg mqtt.Message) {
	req := types.NewRequest().
		SetMetadata(c.createMetadataString(msg)).
		SetData(msg.Payload())
	if c.opts.dynamicMapping {
		req.SetChannel(msg.Topic())
	}
	_, err := c.target.Do(ctx, req)
	if err != nil {
		c.log.Errorf("error processing mqtt message %d , %s", msg.MessageID(), err.Error())
	}
}
func (c *Client) Stop() error {
	c.client.Disconnect(250)
	return nil
}
