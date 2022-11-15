package kafka

import (
	"fmt"
	"strings"

	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-io/kubemq-sources/config"
)

type options struct {
	brokers          []string
	topics           []string
	consumerGroup    string
	saslUsername     string
	saslPassword     string
	saslMechanism    string
	securityProtocol string
	cacert           string
	clientCert       string
	clientKey        string
	insecure         bool
	dynamicMapping   bool
}

func parseOptions(cfg config.Spec) (options, error) {
	m := options{}
	var err error

	m.consumerGroup, err = cfg.Properties.MustParseString("consumer_group")
	if err != nil {
		return m, err
	}
	m.brokers, err = cfg.Properties.MustParseStringList("brokers")
	if err != nil {
		return m, err
	}
	m.topics, err = cfg.Properties.MustParseStringList("topics")
	if err != nil {
		return m, err
	}
	m.saslUsername = cfg.Properties.ParseString("saslUsername", "")
	m.saslPassword = cfg.Properties.ParseString("saslPassword", "")
	m.saslMechanism = cfg.Properties.ParseString("saslMechanism", "")
	m.securityProtocol = cfg.Properties.ParseString("securityProtocol", "")
	m.cacert = cfg.Properties.ParseString("ca_cert", "")
	m.clientCert = cfg.Properties.ParseString("client_certificate", "")
	m.clientKey = cfg.Properties.ParseString("client_key", "")
	m.insecure = cfg.Properties.ParseBool("insecure", false)
	m.dynamicMapping, err = cfg.Properties.MustParseBool("dynamic_mapping")
	if err != nil {
		return options{}, fmt.Errorf("error parsing dynamic_mapping, %w", err)
	}
	return m, nil
}

func (m *options) parseASLMechanism() kafka.SASLMechanism {
	switch strings.ToLower(m.saslMechanism) {
	case "plain":
		return kafka.SASLTypePlaintext
	case "scram-sha-256":
		return kafka.SASLTypeSCRAMSHA256
	case "scram-sha-512":
		return kafka.SASLTypeSCRAMSHA512
	case "gssapi", "gss-api", "gss_api":
		return kafka.SASLTypeGSSAPI
	case "oauth", "0auth bearer":
		return kafka.SASLTypeOAuth
	case "external", "ext":
		return kafka.SASLExtKeyAuth
	default:
		return kafka.SASLTypePlaintext
	}
}

func (m *options) parseSecurityProtocol() (bool, bool) {
	switch strings.ToLower(m.securityProtocol) {
	case "plaintext":
		return false, false
	case "ssl":
		return true, false
	case "sasl_plaintext":
		return false, true
	case "sasl_ssl":
		return true, true
	default:
		return false, false
	}
}
