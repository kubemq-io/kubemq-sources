package kinesis

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"time"
)

const (
	defaultToken           = ""
	defaultSequence        = ""
	DefaultWaitBetweenPull = 5
)

type options struct {
	awsKey            string
	awsSecretKey      string
	region            string
	token             string
	streamARN         string
	consumerARN       string
	sequenceNumber    string
	shardIteratorType string
	shardID           string
	pullDelay         time.Duration
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error

	o.token = cfg.ParseString("token", defaultToken)
	o.sequenceNumber = cfg.ParseString("sequence_number", defaultSequence)
	o.pullDelay = time.Duration(cfg.ParseInt("pull_delay", DefaultWaitBetweenPull))

	o.awsKey, err = cfg.MustParseString("aws_key")
	if err != nil {
		return options{}, fmt.Errorf("error aws_key , %w", err)
	}

	o.awsSecretKey, err = cfg.MustParseString("aws_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error aws_secret_key , %w", err)
	}

	o.region, err = cfg.MustParseString("region")
	if err != nil {
		return options{}, fmt.Errorf("error region , %w", err)
	}
	o.streamARN, err = cfg.MustParseString("stream_arn")
	if err != nil {
		return options{}, fmt.Errorf("error parsing stream_arn, %w", err)
	}
	o.consumerARN, err = cfg.MustParseString("consumer_arn")
	if err != nil {
		return options{}, fmt.Errorf("error parsing consumer_arn, %w", err)
	}

	o.shardIteratorType, err = cfg.MustParseString("shard_iterator_type")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shard_iterator_type, %w", err)
	}
	o.shardID, err = cfg.MustParseString("shard_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shard_id, %w", err)
	}

	return o, nil
}
