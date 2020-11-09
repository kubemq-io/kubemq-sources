package amazonmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

type options struct {
	host           string
	username       string
	password       string
	destination    string
	dynamicMapping bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username, err = cfg.Properties.MustParseString("username")
	if err != nil {
		return options{}, fmt.Errorf("error parsing username , %w", err)
	}
	o.password, err = cfg.Properties.MustParseString("password")
	if err != nil {
		return options{}, fmt.Errorf("error parsing password , %w", err)
	}
	o.destination, err = cfg.Properties.MustParseString("destination")
	if err != nil {
		return options{}, fmt.Errorf("error parsing destination name, %w", err)
	}
	o.dynamicMapping, err = cfg.Properties.MustParseBool("dynamic_mapping")
	if err != nil {
		return options{}, fmt.Errorf("error parsing dynamic_mapping, %w", err)
	}
	return o, nil
}
