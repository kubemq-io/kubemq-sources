package http

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

const (
	defaultHost = ""
	defaultPort = 8080
)

type options struct {
	host string
	port int
	path string
}

func parseOptions(cfg config.Spec) (options, error) {
	m := options{}
	var err error
	m.host = cfg.Properties.ParseString("host", defaultHost)

	m.port, err = cfg.Properties.ParseIntWithRange("port", defaultPort, 1, 65535)
	if err != nil {
		return m, fmt.Errorf("error parsing port value, %w", err)
	}

	m.path = cfg.Properties.ParseString("path", "/")

	return m, nil
}
