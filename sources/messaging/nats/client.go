package nats

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/types"
	"github.com/nats-io/nats.go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	opts   options
	client *nats.Conn
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
	o := setOptions(c.opts.certFile, c.opts.certKey, c.opts.username, c.opts.password, c.opts.token, c.opts.tls, c.opts.timeout)
	c.client, err = nats.Connect(c.opts.url, o)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	} else {
		c.target = target
	}

	_, err := c.client.Subscribe(c.opts.subject, func(m *nats.Msg) {
		go c.processIncomingMessages(ctx, m)
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createMetadataString(msg *nats.Msg) string {
	md := map[string]string{}
	md["subject"] = msg.Subject
	md["replay"] = msg.Reply
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing nats metadata, %s", err.Error())
	}
	return str
}
func (c *Client) processIncomingMessages(ctx context.Context, msg *nats.Msg) {
	req := types.NewRequest().
		SetMetadata(c.createMetadataString(msg)).
		SetData(msg.Data)
	if c.opts.dynamicMapping {
		req.SetChannel(msg.Subject)
	}
	_, err := c.target.Do(ctx, req)
	if err != nil {
		c.log.Errorf("error processing nats message from :%s , %s", msg.Subject, err.Error())
	}
}

func (c *Client) Stop() error {
	c.client.Close()
	return nil
}

func setOptions(sslcertificatefile string, sslcertificatekey string, username string, password string, token string, useTls bool, timeout int) nats.Option {
	return func(o *nats.Options) error {
		if useTls {
			if sslcertificatefile != "" && sslcertificatekey != "" {
				cert, err := tls.X509KeyPair([]byte(sslcertificatefile), []byte(sslcertificatekey))
				if err != nil {
					return fmt.Errorf("nats: error parsing client certificate: %v", err)
				}
				cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
				if err != nil {
					return fmt.Errorf("nats: error parsing client certificate: %v", err)
				}
				if o.TLSConfig == nil {
					o.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
				}
				o.TLSConfig.Certificates = []tls.Certificate{cert}
				o.Secure = true
			} else {
				return errors.New("when using tls make sure to pass file and key")
			}
		}
		if username != "" {
			o.User = username
		}
		if password != "" {
			o.Password = password
		}
		if token != "" {
			o.Token = token
		}
		if timeout != 0 {
			o.Timeout = time.Duration(timeout) * time.Second
		}

		return nil
	}
}
