package rabbitmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/nats-io/nuid"
)

type options struct {
	url            string
	queue          string
	consumer       string
	requeueOnError bool
	autoAck        bool
	exclusive      bool
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.url, err = cfg.Properties.MustParseString("url")
	if err != nil {
		return options{}, fmt.Errorf("error parsing url, %w", err)
	}
	o.queue, err = cfg.Properties.MustParseString("queue")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue name, %w", err)
	}
	o.consumer = cfg.Properties.ParseString("consumer", nuid.Next())
	o.requeueOnError = cfg.Properties.ParseBool("requeue_on_error", false)
	o.autoAck = cfg.Properties.ParseBool("auto_ack", false)
	o.exclusive = cfg.Properties.ParseBool("exclusive", false)
	return o, nil
}
