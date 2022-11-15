package rabbitmq

import (
	"fmt"

	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/nats-io/nuid"
)

type options struct {
	url               string
	dynamicMapping    bool
	queue             string
	consumer          string
	requeueOnError    bool
	autoAck           bool
	exclusive         bool
	clientCertificate string
	clientKey         string
	caCert            string
	insecure          bool
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
	o.dynamicMapping = cfg.Properties.ParseBool("dynamic_mapping", false)

	o.consumer = cfg.Properties.ParseString("consumer", nuid.Next())
	o.requeueOnError = cfg.Properties.ParseBool("requeue_on_error", false)
	o.autoAck = cfg.Properties.ParseBool("auto_ack", false)
	o.exclusive = cfg.Properties.ParseBool("exclusive", false)
	o.caCert = cfg.Properties.ParseString("ca_cert", "")
	o.clientCertificate = cfg.Properties.ParseString("client_certificate", "")
	o.clientKey = cfg.Properties.ParseString("client_key", "")
	o.insecure = cfg.Properties.ParseBool("insecure", false)
	return o, nil
}
