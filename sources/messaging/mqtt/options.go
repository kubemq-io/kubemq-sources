package mqtt

import (
	"fmt"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/nats-io/nuid"
)

type options struct {
	host           string
	dynamicMapping bool
	topic          string
	username       string
	password       string
	clientId       string
	qos            int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.topic, err = cfg.Properties.MustParseString("topic")
	if err != nil {
		return options{}, fmt.Errorf("error parsing topic, %w", err)
	}
	o.dynamicMapping, err = cfg.Properties.MustParseBool("dynamic_mapping")
	if err != nil {
		return options{}, fmt.Errorf("error parsing dynamic mapping, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.clientId = cfg.Properties.ParseString("client_id", nuid.Next())
	o.qos, err = cfg.Properties.ParseIntWithRange("qos", 0, 0, 2)
	if err != nil {
		return options{}, fmt.Errorf("error parsing qos, %w", err)
	}

	return o, nil
}
