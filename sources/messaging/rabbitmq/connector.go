package rabbitmq

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.rabbitmq").
		SetDescription("RabbitMQ source properties").
		SetName("RabbitMQ").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("queue", "pub/sub").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("url").
				SetTitle("URL Address").
				SetDescription("Set RabbitMQ connection string").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("queue").
				SetTitle("Queue Name").
				SetDescription("Set queue name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("dynamic_mapping").
				SetDescription("Set Queue/Channel dynamic mapping").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("consumer").
				SetTitle("Consumer ID").
				SetDescription("Set consumer tag value").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("requeue_on_error").
				SetTitle("Requeue on Error").
				SetDescription("Set requeue message on error").
				SetMust(false).
				SetDefault("false"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("auto_ack").
				SetTitle("Automatic Ack").
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
