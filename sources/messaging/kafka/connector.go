package kafka

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.kafka").
		SetDescription("Kafka source properties").
		SetName("Kafka").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("pub/sub", "streaming").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("brokers").
				SetTitle("Brokers Address").
				SetDescription("Set Brokers list").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("topics").
				SetDescription("Set Topics list").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("consumer_group").
				SetDescription("Set Consumer Group").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_username").
				SetTitle("SASL Username").
				SetDescription("Set SASL Username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_password").
				SetTitle("SASL Password").
				SetDescription("Set SASL Password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_mechanism").
				SetTitle("SASL Mechanism").
				SetDescription("SCRAM-SHA-256, SCRAM-SHA-512, plain, 0Auth bearer, or GSS-API").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("security_protocol").
				SetTitle("Security Protocol").
				SetDescription("plaintext, SASL-plaintext, SASL-SSL, SSL").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("ca_cert").
				SetDescription("Set TLS CA Certificate").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("client_certificate").
				SetDescription("Set TLS Client PEM data").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("client_key").
				SetDescription("Set TLS Client Key PEM data").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("insecure").
				SetDescription("Set self-signed SSL Certificate").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("dynamic_mapping").
				SetDescription("Set Topic/Channel dynamic mapping").
				SetMust(true).
				SetDefault("true"),
		)
}
