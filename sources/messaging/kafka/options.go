package kafka

import (
	"github.com/kubemq-hub/kubemq-sources/config"
)

type options struct {
	brokers       []string
	topics        []string
	consumerGroup string
	saslUsername  string
	saslPassword  string
}

func parseOptions(cfg config.Spec) (options, error) {
	m := options{}
	var err error

	m.consumerGroup, err = cfg.Properties.MustParseString("consumer_group")
	if err != nil {
		return m, err
	}
	m.brokers, err = cfg.Properties.MustParseStringList("brokers")
	if err != nil {
		return m, err
	}
	m.topics, err = cfg.Properties.MustParseStringList("topics")
	if err != nil {
		return m, err
	}
	m.saslUsername = cfg.Properties.ParseString("saslUsername", "")
	m.saslPassword = cfg.Properties.ParseString("saslPassword", "")

	return m, nil
}
