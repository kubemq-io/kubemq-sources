package amazonmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"time"
)

const (
	defaultSubTimeout = 5
)


type options struct {
	host     string
	username string
	password string
	subTimeout  time.Duration
	destination string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.Properties.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username ,err= cfg.Properties.MustParseString("username")
	if err != nil {
		return options{}, fmt.Errorf("error parsing username , %w", err)
	}
	o.password ,err= cfg.Properties.MustParseString("password")
	if err != nil {
		return options{}, fmt.Errorf("error parsing password , %w", err)
	}
	o.destination, err = cfg.Properties.MustParseString("destination")
	if err != nil {
		return options{}, fmt.Errorf("error parsing destination name, %w", err)
	}
	o.subTimeout = time.Duration(cfg.Properties.ParseInt("sub_timeout", defaultSubTimeout))
	return o, nil
}
