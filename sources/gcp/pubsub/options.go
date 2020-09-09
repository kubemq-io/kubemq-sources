package pubsub

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

type options struct {
	projectID    string
	credentials  string
	subscriberID string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.projectID, err = cfg.MustParseString("project_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.subscriberID, err = cfg.MustParseString("subscriber_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing project_id, %w", err)
	}
	o.credentials, err = cfg.MustParseString("credentials")
	if err != nil {
		return options{}, err
	}

	return o, nil
}
