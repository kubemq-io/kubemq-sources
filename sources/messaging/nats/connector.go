package nats

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.nats").
		SetDescription("nats source properties").
		SetName("NATS").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("queue","pub/sub").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Set nats url connection").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("subject").
				SetDescription("Set subject").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("tls").
				SetDescription("Set if use tls").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("dynamic_mapping").
				SetDescription("Set Subject/Channel dynamic mapping").
				SetMust(true).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("tls").
				SetOptions([]string{"true", "false"}).
				SetDescription("Set tls conditions").
				SetMust(true).
				SetDefault("false").
				NewCondition("true", []*common.Property{
					common.NewProperty().
						SetKind("multilines").
						SetName("cert_key").
						SetDescription("Set certificate key").
						SetMust(false).
						SetDefault(""),
					common.NewProperty().
						SetKind("multilines").
						SetName("cert_file").
						SetDescription("Set certificate file").
						SetMust(false).
						SetDefault(""),
				}),
		)
}
