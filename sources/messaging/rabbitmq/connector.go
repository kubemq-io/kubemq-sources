package rabbitmq

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("source.messaging.rabbitmq").
		SetDescription("RabbitMQ source properties").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetDescription("Set RabbitMQ connection string").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("queue").
				SetDescription("Set queue name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("consumer").
				SetDescription("Set consumer tag value").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("requeue_on_error").
				SetDescription("Set requeue message on error").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("auto_ack").
				SetDescription("Set auto ack upon receive message").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("exclusive").
				SetDescription("Set exclusive subscription").
				SetMust(false).
				SetDefault("false"),
		)
}
