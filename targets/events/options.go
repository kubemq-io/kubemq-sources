package events

import (
	"fmt"

	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/nats-io/nuid"
)

const (
	defaultAddress = "localhost:50000"
)

type options struct {
	host      string
	port      int
	clientId  string
	authToken string
	channel   string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, o.port, err = cfg.Properties.MustParseAddress("address", defaultAddress)
	if err != nil {
		return options{}, fmt.Errorf("error parsing address value, %w", err)
	}
	o.channel, err = cfg.Properties.MustParseString("channel")
	if err != nil {
		return options{}, fmt.Errorf("error parsing channel value, %w", err)
	}
	o.authToken = cfg.Properties.ParseString("auth_token", "")
	o.clientId = cfg.Properties.ParseString("client_id", nuid.Next())
	return o, nil
}
