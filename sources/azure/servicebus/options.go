package servicebus

import (
	"fmt"

	"github.com/kubemq-io/kubemq-sources/config"
)

type options struct {
	connectionString string
	queueName        string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	endPoint, err := cfg.Properties.MustParseString("end_point")
	if err != nil {
		return options{}, fmt.Errorf("error parsing end_point , %w", err)
	}
	sharedAccessKeyName, err := cfg.Properties.MustParseString("shared_access_key_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shared_access_key_name , %w", err)
	}
	sharedAccessKey, err := cfg.Properties.MustParseString("shared_access_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shared_access_key , %w", err)
	}
	o.queueName, err = cfg.Properties.MustParseString("queue_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue_name , %w", err)
	}
	o.connectionString = fmt.Sprintf("Endpoint=%s;SharedAccessKeyName=%s;SharedAccessKey=%s", endPoint, sharedAccessKeyName, sharedAccessKey)
	return o, nil
}
