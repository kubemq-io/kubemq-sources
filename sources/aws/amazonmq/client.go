package amazonmq

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-stomp/stomp"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	opts   options
	conn   *stomp.Conn
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

	netConn, err := tls.Dial("tcp", c.opts.host, &tls.Config{})
	if err != nil {
		return err
	}

	c.conn, err = stomp.Connect(netConn, stomp.ConnOpt.Login(c.opts.username, c.opts.password))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) createMetadataString(msg *stomp.Message) string {
	md := map[string]string{}
	md["destination"] = msg.Destination
	md["content_type"] = msg.ContentType
	if msg.Err != nil {
		md["error"] = msg.Err.Error()
	}
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing stomp metadata, %s", err.Error())
	}
	return str
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	} else {
		c.target = target
	}
	errCh := make(chan error, 1)
	go func() {
		subscription, err := c.conn.Subscribe(c.opts.destination, stomp.AckAuto)
		if err != nil {
			errCh <- fmt.Errorf("error subscription to amazonmq destination, %w", err)
			return
		}
		defer func() {
			_ = subscription.Unsubscribe()
		}()
		for {
			select {
			case msg := <-subscription.C:
				req := types.NewRequest().SetMetadata(c.createMetadataString(msg)).SetData(msg.Body)
				if c.opts.dynamicMapping {
					req.SetChannel(msg.Destination)
				}
				_, err := c.target.Do(ctx, req)
				if err != nil {
					c.log.Errorf("error processing amazonmq message, %s", err.Error())
				}
			case err := <-errCh:
				c.log.Errorf("error on amazonmq connection, %v", err)
			case <-ctx.Done():
				return
			}
		}

	}()

	return nil
}

func (c *Client) Stop() error {
	return c.conn.Disconnect()
}
