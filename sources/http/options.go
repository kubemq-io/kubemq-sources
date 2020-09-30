package http

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/pkg/builder"
)

type options struct {
	methods []string
	path    string
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
	return o, nil
}

func Connector() *builder.Connector {
	return builder.NewConnector().
		SetKind("source.http").
		SetDescription("http source properties").
		AddProperty(
			builder.NewProperty().
				SetKind("string").
				SetName("methods").
				SetDescription("list of supported methods separated by a comma").
				SetMust(true),
		).
		AddProperty(
			builder.NewProperty().
				SetKind("string").
				SetName("path").
				SetDescription("http endpoint path").
				SetMust(true).
				SetDefault("/"),
		)
}
