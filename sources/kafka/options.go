package kafka

import (
	"github.com/kubemq-hub/kubemq-source-connectors/config"
)

type options struct {
	brokers       []string
	topics        []string
	consumerGroup string
	saslUsername  string
	saslPassword  string
}

func parseOptions(cfg config.Metadata) (options, error) {
	m := options{}
	var err error

	m.consumerGroup, err = cfg.MustParseString("consumerGroup")
	if err != nil {
		return m, err
	}
	m.brokers, err = cfg.MustParseStringList("brokers")
	if err != nil {
		return m, err
	}
	m.topics, err = cfg.MustParseStringList("topics")
	if err != nil {
		return m, err
	}
	m.saslUsername = cfg.ParseString("saslUsername", "")
	m.saslPassword = cfg.ParseString("saslPassword", "")

	return m, nil
}
