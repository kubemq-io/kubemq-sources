package mqtt

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/nats-io/nuid"
)

type options struct {
	host     string
	topic    string
	username string
	password string
	clientId string
	qos      int
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.topic, err = cfg.MustParseString("topic")
	if err != nil {
		return options{}, fmt.Errorf("error parsing topic, %w", err)
	}
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")
	o.clientId = cfg.ParseString("client_id", nuid.Next())
	o.qos, err = cfg.ParseIntWithRange("qos", 0, 0, 2)
	if err != nil {
		return options{}, fmt.Errorf("error parsing qos, %w", err)
	}

	return o, nil
}
