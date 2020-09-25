package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/kubemq-hub/kubemq-sources/types"
	"time"
)

type Client struct {
	name   string
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

func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name

	c.log = logger.NewLogger(c.name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.opts.region),
		Credentials: credentials.NewStaticCredentials(c.opts.sqsKey, c.opts.sqsSecretKey, c.opts.token),
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
							req := types.NewRequest().SetData(b)
							_, err := target.Do(ctx, req)
							if err != nil {
								c.log.Errorf("error processing request %s", err.Error())
							}
						}
					}
				}
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
