package msk

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.msk").
		SetDescription("AWS MSK source properties").
		SetName("MSK").
		SetProvider("AWS").
		SetCategory("Messaging").
		SetTags("streaming", "kafka", "cloud", "managed").
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
				SetTitle("Consumer Group").
				SetDescription("Set Consumer Group").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_username").
				SetTitle("SASL Username").
				SetDescription("Set Saal Username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sasl_password").
				SetTitle("SASL Password").
				SetDescription("Set Saal Password").
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
