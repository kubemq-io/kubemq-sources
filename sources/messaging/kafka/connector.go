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
		SetTags("pub/sub","streaming").
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
			SetKind("bool").
			SetName("dynamic_mapping").
			SetDescription("Set Topic/Channel dynamic mapping").
			SetMust(true).
			SetDefault("true"),
	)
}
