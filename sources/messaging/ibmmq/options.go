package ibmmq

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"math"
)

const (
	defaultCertificateLabel = ""
	defaultPassword         = ""
	defaultKeyRepository    = ""
	defaultWaitPullTime     = 100
	defaultPort             = 1414
	defaultTlsClientAuth    = "NONE"
)

type options struct {
	qMName           string
	hostname         string
	portNumber       int
	channelName      string
	userName         string
	keyRepository    string
	certificateLabel string
	queueName        string
	Password         string
	transportType    int
	pullDelay       int32
	tlsClientAuth    string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.qMName, err = cfg.Properties.MustParseString("queue_manager_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue_manager_name, %w", err)
	}
	o.hostname, err = cfg.Properties.MustParseString("host_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host_name, %w", err)
	}
	o.channelName, err = cfg.Properties.MustParseString("channel_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing channel_name, %w", err)
	}
	o.userName, err = cfg.Properties.MustParseString("username")
	if err != nil {
		return options{}, fmt.Errorf("error parsing username, %w", err)
	}
	o.queueName, err = cfg.Properties.MustParseString("queue_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing queue_name, %w", err)
	}
	o.certificateLabel = cfg.Properties.ParseString("certificate_label", defaultCertificateLabel)
	timeToLive, err := cfg.Properties.ParseIntWithRange("wait_pull_time", defaultWaitPullTime, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error parsing wait_pull_time, %w", err)
	}
	o.pullDelay = int32(timeToLive)
	o.portNumber = cfg.Properties.ParseInt("port_number", defaultPort)
	o.Password = cfg.Properties.ParseString("password", defaultPassword)
	o.tlsClientAuth = cfg.Properties.ParseString("tls_client_auth", defaultTlsClientAuth)
	o.keyRepository = cfg.Properties.ParseString("key_repository", defaultKeyRepository)
	return o, nil
}
