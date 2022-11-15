package pubsub

import (
	"fmt"

	"github.com/kubemq-io/kubemq-sources/config"
)

type options struct {
	projectID    string
	credentials  string
	subscriberID string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.projectID, err = cfg.Properties.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.subscriberID, err = cfg.Properties.MustParseString("subscriber_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.credentials, err = cfg.Properties.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}

	return o, nil
}
