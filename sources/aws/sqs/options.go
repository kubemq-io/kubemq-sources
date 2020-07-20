package sqs

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

const (
	DefaultMaxReceive = 0
	DefaultMaxNumberOfMessages = 0
	DefaultWaitTimeSeconds = 0
	DefaultVisibility = 0
	DefaultToken = ""
)

type options struct {
	sqsKey          string
	sqsSecretKey    string
	region          string
	maxReceiveCount int
	concurrency     int
	queue           string
	deadLetterQueue string
	token string
	visibility int64
	maxNumberOfMessages int64
	waitTimeSeconds int64
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.sqsKey, err = cfg.MustParseString("sqs_key")
	if err != nil {
		return options{}, fmt.Errorf("error sqsKey , %w", err)
	}

	o.sqsSecretKey, err = cfg.MustParseString("sqs_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error sqsSecretKey , %w", err)
	}


	o.region, err = cfg.MustParseString("region")
	if err != nil {
		return options{}, fmt.Errorf("error region , %w", err)
	}

	o.maxReceiveCount = cfg.ParseInt("max_receive", DefaultMaxReceive)
	o.deadLetterQueue = cfg.ParseString("dead_letter", "")

	o.visibility = int64(cfg.ParseInt("visibility", DefaultVisibility))
	o.maxNumberOfMessages = int64(cfg.ParseInt("max_number_of_messages", DefaultMaxNumberOfMessages))
	o.waitTimeSeconds = int64(cfg.ParseInt("wait_time_seconds", DefaultWaitTimeSeconds))
	
	o.concurrency, err = cfg.MustParseIntWithRange("concurrency", 1, 100)
	if err != nil {
		return options{}, fmt.Errorf("error parsing concurrency value, %w", err)
	}
	o.token = cfg.ParseString("token", DefaultToken)
	o.queue, err = cfg.MustParseString("queue")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue, %w", err)
	}
	
	return o, nil
}
