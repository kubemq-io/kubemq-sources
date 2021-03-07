package activemq

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.activemq").
		SetDescription("ActiveMQ source properties").
		SetName("ActiveMQ").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("queue","streaming").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetTitle("Host Address").
				SetDescription("Set ActiveMQ Host connection").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("destination").
				SetDescription("Set Destination").
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
		)
}
