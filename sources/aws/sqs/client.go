package sqs

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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
	log    *logger.Logger
	client *sqs.SQS
	ctx    context.Context
	cancel context.CancelFunc
	target middleware.Middleware
}

func New() *Client {
	return &Client{}
}

func (c *Client) createMetadataString(msg *sqs.Message) string {
	md := map[string]string{}
	if msg.MessageId != nil {
		md["message_id"] = *msg.MessageId
	}
	if len(msg.MessageAttributes) > 0 {
		ma, err := json.Marshal(msg.MessageAttributes)
		if err != nil {
			return fmt.Sprintf("error parsing MessageAttributes, %s", err.Error())
		}
		md["message_attributes"] = string(ma)
	}
	if len(msg.Attributes) > 0 {
		a, err := json.Marshal(msg.Attributes)
		if err != nil {
			return fmt.Sprintf("error parsing Attributes, %s", err.Error())
		}
		md["attributes"] = string(a)
	}
	md["receipt_handler"] = *msg.ReceiptHandle
	str, err := json.MarshalToString(md)
	if err != nil {
		return fmt.Sprintf("error parsing stomp metadata, %s", err.Error())
	}
	return str
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
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.opts.region),
		Credentials: credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}

	svc := sqs.New(sess)
	c.client = svc
	c.ctx, c.cancel = context.WithCancel(ctx)
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	if target == nil {
		return fmt.Errorf("invalid target received, cannot be nil")
	} else {
		c.target = target
	}

	r := sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(c.opts.maxNumberOfMessages),
		QueueUrl:            aws.String(c.opts.queue),
		WaitTimeSeconds:     aws.Int64(c.opts.waitTimeSeconds),
		VisibilityTimeout:   aws.Int64(c.opts.visibilityTimeout),
	}
	go func() {
		for {
			select {
			case <-time.After(c.opts.pullDelay * time.Millisecond):
				r, err := c.client.ReceiveMessageWithContext(ctx, &r)
				if err != nil {
					c.log.Errorf("error receiving request %v", err)
				} else {
					for _, message := range r.Messages {
						b, err := json.Marshal(message)
						if err != nil {
							c.log.Errorf("failed to parse message on error %s", err.Error())
						} else {
							req := types.NewRequest().SetMetadata(c.createMetadataString(message)).SetData(b)
							_, err := target.Do(ctx, req)
							if err != nil {
								c.log.Errorf("error processing request %s", err.Error())
							} else {
								_, err = c.client.DeleteMessage(&sqs.DeleteMessageInput{
									QueueUrl:      aws.String(c.opts.queue),
									ReceiptHandle: message.ReceiptHandle,
								})
								c.log.Errorf("failed to delete message on error %s", err.Error())
							}
						}
					}
				}
			case <-ctx.Done():
				return
			case <-c.ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (c *Client) Stop() error {
	c.cancel()
	return nil
}
