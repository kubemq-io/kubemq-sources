package kinesis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/kubemq-hub/builder/connector/common"
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
	client *kinesis.Kinesis
	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Client {
	return &Client{}
}
func (c *Client) Connector() *common.Connector {
	return Connector()
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
		Credentials: credentials.NewStaticCredentials(c.opts.awsKey, c.opts.awsSecretKey, c.opts.token),
	})
	if err != nil {
		return err
	}

	svc := kinesis.New(sess)
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.client = svc
	return nil
}
func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {

	sp := &kinesis.StartingPosition{
		Type: aws.String(c.opts.ShardIteratorType),
	}
	if c.opts.sequenceNumber != "" {
		sp.SequenceNumber = aws.String(c.opts.sequenceNumber)
	}
	si := &kinesis.SubscribeToShardInput{
		ConsumerARN:      aws.String(c.opts.consumerARN),
		ShardId:          aws.String(c.opts.shardID),
		StartingPosition: sp,
	}

	response, err := c.client.SubscribeToShardWithContext(ctx, si)
	if err != nil {
		return fmt.Errorf("failed to connect to subscribe to shard on err :%s", err.Error())
	}

	go func() {
		for {
			select {
			case <-time.After(time.Duration(c.opts.pullDelay) * time.Millisecond):
				for events := range response.GetStream().Events() {
					if events == nil {
						c.log.Errorf("failed to receive events on error %s", err.Error())
					}
					b, err := json.Marshal(events)
					if err != nil {
						c.log.Errorf("failed to parse record on error %s", err.Error())
					} else {
						req := types.NewRequest().SetData(b)
						_, err := target.Do(ctx, req)
						if err != nil {
							c.log.Errorf("error processing request %s", err.Error())
						}
					}
				}
			case <-c.ctx.Done():
				response.GetStream().Close()
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
