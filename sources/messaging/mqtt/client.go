package mqtt

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	defaultConnectTimeout = 5 * time.Second
)

type Client struct {
	opts           options
	client         mqtt.Client
	log            *logger.Logger
	target         middleware.Middleware
	isConnected    *atomic.Bool
	reconnectCount *atomic.Int32
}

func New() *Client {
	return &Client{
		isConnected:    atomic.NewBool(false),
		reconnectCount: atomic.NewInt32(0),
	}
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
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(1 * time.Second)
	opts.SetMaxReconnectInterval(24 * time.Hour)
	opts.SetConnectRetry(true)
	opts.SetOnConnectHandler(c.onConnect)
	opts.SetConnectionLostHandler(c.onConnectionLost)
	opts.SetReconnectingHandler(c.onReconnectingHandler)
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error connecting to mqtt broker, %w", token.Error())
	}
	c.log.Infof("connected to mqtt broker at %s", c.opts.host)
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	c.log.Infof("client started")
	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	} else {
		c.target = target
	}

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
	c.log.Infof("sending request to target %s", req.String())
	resp, err := c.target.Do(ctx, req)
	if err != nil {
		c.log.Errorf("error processing mqtt message %d , %s", msg.MessageID(), err.Error())
	} else {
		c.log.Infof("target response %s", resp.String())
	}
	if resp != nil {
		c.log.Infof("target response %s", resp.String())
	}
}

func (c *Client) Stop() error {
	c.log.Infof("client stopped")
	c.client.Disconnect(250)
	return nil
}

func (c *Client) onConnectionLost(client mqtt.Client, err error) {
	c.log.Errorf("mqtt client connection lost, error: %s", err.Error())
	c.isConnected.Store(false)
}

func (c *Client) onConnect(client mqtt.Client) {
	c.log.Infof("mqtt client connected")
	c.isConnected.Store(true)
	c.log.Infof("subscribing to topic %s", c.opts.topic)
	c.client.Subscribe(c.opts.topic, byte(c.opts.qos), func(client mqtt.Client, message mqtt.Message) {
		c.log.Infof("received message from topic %s, payload %s", message.Topic(), string(message.Payload()))
		ctx := context.Background()
		go c.processIncomingMessages(ctx, message)
	})
	c.log.Infof("subscribed to topic %s", c.opts.topic)
	c.reconnectCount.Store(0)
}

func (c *Client) onReconnectingHandler(client mqtt.Client, opts *mqtt.ClientOptions) {
	c.reconnectCount.Inc()
	c.log.Warnf("mqtt client reconnecting to broker, attempt: %d", c.reconnectCount.Load())

}
