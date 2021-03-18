package http

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

type options struct {
	methods        []string
	path           string
	dynamicMapping bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.methods, err = cfg.Properties.MustParseStringList("methods")
	if err != nil {
		return options{}, fmt.Errorf("error parsing methods list value, %w", err)
	}

	o.path, err = cfg.Properties.MustParseString("path")
	if err != nil {
		return options{}, fmt.Errorf("error parsing path value, %w", err)
	}
	o.dynamicMapping = cfg.Properties.ParseBool("dynamic_mapping",false)
	return o, nil
}
