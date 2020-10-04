package sqs

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"time"
)

const (
	DefaultMaxNumberOfMessages = 1
	DefaultWaitTimeSeconds     = 0
	DefaultVisibilityTimeout   = 0
	DefaultToken               = ""
	DefaultWaitBetweenPull     = 5
)

type options struct {
	awsKey              string
	awsSecretKey        string
	region              string
	queue               string
	token               string
	visibilityTimeout   int64
	maxNumberOfMessages int64
	waitTimeSeconds     int64
	pullDelay           time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.awsKey, err = cfg.Properties.MustParseString("aws_key")
	if err != nil {
		return options{}, fmt.Errorf("error aws_key , %w", err)
	}

	o.awsSecretKey, err = cfg.Properties.MustParseString("aws_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error aws_secret_key , %w", err)
	}

	o.region, err = cfg.Properties.MustParseString("region")
	if err != nil {
		return options{}, fmt.Errorf("error region , %w", err)
	}

	o.visibilityTimeout = int64(cfg.Properties.ParseInt("visibility_timeout", DefaultVisibilityTimeout))
	o.maxNumberOfMessages = int64(cfg.Properties.ParseInt("max_number_of_messages", DefaultMaxNumberOfMessages))
	o.waitTimeSeconds = int64(cfg.Properties.ParseInt("wait_time_seconds", DefaultWaitTimeSeconds))

	o.token = cfg.Properties.ParseString("token", DefaultToken)
	o.queue, err = cfg.Properties.MustParseString("queue")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue, %w", err)
	}
	o.pullDelay = time.Duration(cfg.Properties.ParseInt("pull_delay", DefaultWaitBetweenPull))
	return o, nil
}
