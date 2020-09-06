package sqs

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"time"
)

const (
	DefaultMaxNumberOfMessages = 1
	DefaultWaitTimeSeconds     = 0
	DefaultVisibilityTimeout          = 0
	DefaultToken               = ""
	DefaultWaitBetweenPull     = 60
)

type options struct {
	sqsKey       string
	sqsSecretKey string
	region       string

	queue               string
	token               string
	visibilityTimeout          int64
	maxNumberOfMessages int64
	waitTimeSeconds     int64
	pullDelay           time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.sqsKey, err = cfg.MustParseString("aws_key")
	if err != nil {
		return options{}, fmt.Errorf("error aws_key , %w", err)
	}

	o.sqsSecretKey, err = cfg.MustParseString("aws_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error aws_secret_key , %w", err)
	}

	o.region, err = cfg.MustParseString("region")
	if err != nil {
		return options{}, fmt.Errorf("error region , %w", err)
	}


	o.visibilityTimeout = int64(cfg.ParseInt("visibility_timeout", DefaultVisibilityTimeout))
	o.maxNumberOfMessages = int64(cfg.ParseInt("max_number_of_messages", DefaultMaxNumberOfMessages))
	o.waitTimeSeconds = int64(cfg.ParseInt("wait_time_seconds", DefaultWaitTimeSeconds))

	o.token = cfg.ParseString("token", DefaultToken)
	o.queue, err = cfg.MustParseString("queue")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue, %w", err)
	}
	o.pullDelay = time.Duration(cfg.ParseInt("pull_delay", DefaultWaitBetweenPull))
	return o, nil
}
