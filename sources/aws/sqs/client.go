package sqs

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kubemq-hub/kubemq-sources/targets"
	"github.com/kubemq-hub/kubemq-sources/types"
	"strconv"
	
	"errors"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

type Client struct {
	name   string
	opts   options
	client *sqs.SQS
	log    *logger.Logger
	target targets.Target
	doneCh chan struct{}
}

func New() *Client {
	return &Client{}
	
}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.doneCh = make(chan struct{})
	
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.opts.region),
		Credentials: credentials.NewStaticCredentials(c.opts.sqsKey, c.opts.sqsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}
	
	svc := sqs.New(sess)
	
	c.client = svc
	
	return nil
}



func (c *Client) Start(ctx context.Context, target targets.Target) error {
	if target == nil {
		return errInvalidTarget
	} else {
		c.target = target
	}
	for i := 0; i < c.opts.concurrency; i++ {
		go c.run(ctx)
	}
	return nil
}

func (c *Client) run(ctx context.Context) {
	for {
		queueMessages, err := c.getQueueMessages(ctx)
		if err != nil {
			c.log.Error(err.Error())
			return
		}
		for _, message := range queueMessages {
			err := c.processQueueMessage(ctx, message)
			if err != nil {
				c.log.Error(err.Error())
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-c.doneCh:
			return
		default:
		
		}
	}
	
}
func (c *Client) getQueueMessages(ctx context.Context) ([]*sqs.Message, error) {
	
	result, err := c.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.opts.queue,
		MaxNumberOfMessages: aws.Int64(c.opts.maxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(c.opts.visibility),
		WaitTimeSeconds:     aws.Int64(c.opts.waitTimeSeconds),
		MessageAttributeNames: []*string{aws.String("All")},
	})
	if err != nil {
		return nil, err
	}
	
	return result.Messages, nil
}

func (c *Client) processQueueMessage(ctx context.Context, msg *sqs.Message) error {
	req, err := types.ParseRequestFromSQSMessage(msg)
	if err != nil {
		return fmt.Errorf("invalid request format, %w", err)
	}
	_, err = c.target.Do(ctx, req)
	if err != nil {
		return err
	}
	_, err = c.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &c.opts.queue,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		return err
	}
	return nil
	
}

func (c *Client) Stop() error {
	close(c.doneCh)
	return nil
}

func (c *Client) SetQueueAttributes(ctx context.Context, QueueUrl string) error {
	if c.opts.maxReceiveCount > 0 && len(c.opts.deadLetterQueue) > 0 {
		policy := map[string]string{
			"deadLetterTargetArn": c.opts.deadLetterQueue,
			"maxReceiveCount":     strconv.Itoa(c.opts.maxReceiveCount),
		}
		b, err := json.Marshal(policy)
		if err != nil {
			return fmt.Errorf("failed to marshal policy on err :%s", err.Error())
		}
		
		_, err = c.client.SetQueueAttributesWithContext(ctx, &sqs.SetQueueAttributesInput{
			Attributes: map[string]*string{
				sqs.QueueAttributeNameRedrivePolicy: aws.String(string(b)),
			},
			QueueUrl: aws.String(QueueUrl),
		})
		if err != nil {
			return fmt.Errorf("failed to SetQueueAttributesWithContext err :%s", err.Error())
		}
	}
	return fmt.Errorf("failed to SetQueueAttributesWithContext need to verify max_receive and dead_letter exists")
}
