package activemq

import (
	"fmt"
	"github.com/kubemq-io/kubemq-sources/config"
)

type options struct {
	host           string
	destination    string
	username       string
	password       string
	dynamicMapping bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.destination, err = cfg.Properties.MustParseString("destination")
	if err != nil {
		return options{}, fmt.Errorf("error parsing destination, %w", err)
	}
	o.username = cfg.Properties.ParseString("username", "")
	o.password = cfg.Properties.ParseString("password", "")
	o.dynamicMapping, err = cfg.Properties.MustParseBool("dynamic_mapping")
	if err != nil {
		return options{}, fmt.Errorf("error parsing dynamic_mapping, %w", err)
	}
	return o, nil
}
