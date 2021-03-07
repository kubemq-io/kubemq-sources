package amazonmq

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.amazonmq").
		SetDescription("AWS AmazonMQ source properties").
		SetName("AmazonMQ").
		SetProvider("AWS").
		SetCategory("Messaging").
		SetTags("queue","pub/sub","cloud","managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetTitle("Host Address").
				SetDescription("Set AmazonMQ host").
				SetMust(true),

		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set username").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set password").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("destination").
				SetDescription("Set destination").
				SetMust(true).
				SetDefault(""),
		).AddProperty(
		common.NewProperty().
			SetKind("bool").
			SetName("dynamic_mapping").
			SetDescription("Set Topic/Channel dynamic mapping").
			SetMust(true).
			SetDefault("true"),
	)
}
