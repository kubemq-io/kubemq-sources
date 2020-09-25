package events

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/nats-io/nuid"
)

const (
	defaultHost           = "localhost"
	defaultPort           = 50000
	defaultMaxConcurrency = 100
	defaultConcurrency    = 1
)

type options struct {
	host           string
	port           int
	clientId       string
	authToken      string
	concurrency    int
	defaultChannel string
}

func parseOptions(cfg config.Spec) (options, error) {
	m := options{}
	var err error
	m.host = cfg.Properties.ParseString("host", defaultHost)

	m.port, err = cfg.Properties.ParseIntWithRange("port", defaultPort, 1, 65535)
	if err != nil {
		return options{}, fmt.Errorf("error parsing port value, %w", err)
	}
	m.authToken = cfg.Properties.ParseString("auth_token", "")
	m.clientId = cfg.Properties.ParseString("client_id", nuid.Next())
	m.defaultChannel = cfg.Properties.ParseString("default_channel", "")
	m.concurrency, err = cfg.Properties.ParseIntWithRange("concurrency", defaultConcurrency, 1, defaultMaxConcurrency)
	if err != nil {
		return options{}, fmt.Errorf("error parsing concurrency value, %w", err)
	}
	return m, nil
}
