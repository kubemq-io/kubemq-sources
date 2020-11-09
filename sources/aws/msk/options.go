package msk

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

const (
	defaultUsername = ""
	defaultPassword = ""
)

type options struct {
	brokers        []string
	topics         []string
	consumerGroup  string
	saslUsername   string
	saslPassword   string
	dynamicMapping bool
}

func parseOptions(cfg config.Spec) (options, error) {
	m := options{}
	var err error

	m.consumerGroup, err = cfg.Properties.MustParseString("consumer_group")
	if err != nil {
		return options{}, fmt.Errorf("error parsing consumer_group name, %w", err)
	}
	m.brokers, err = cfg.Properties.MustParseStringList("brokers")
	if err != nil {
		return options{}, fmt.Errorf("error parsing brokers name, %w", err)
	}
	m.topics, err = cfg.Properties.MustParseStringList("topics")
	if err != nil {
		return options{}, fmt.Errorf("error parsing topics name, %w", err)
	}
	m.saslUsername = cfg.Properties.ParseString("sasl_username", defaultUsername)
	m.saslPassword = cfg.Properties.ParseString("sasl_password", defaultPassword)
	m.dynamicMapping, err = cfg.Properties.MustParseBool("dynamic_mapping")
	if err != nil {
		return options{}, fmt.Errorf("error parsing dynamic_mapping, %w", err)
	}

	return m, nil
}
