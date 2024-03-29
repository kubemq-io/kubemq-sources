package servicebus

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.servicebus").
		SetDescription("Azure ServiceBus Source").
		SetName("ServiceBus").
		SetProvider("Azure").
		SetCategory("Messaging").
		SetTags("queue", "pub/sub", "cloud", "managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("end_point").
				SetTitle("Endpoint Address").
				SetDescription("Set ServiceBus end point").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key_name").
				SetTitle("Access Key Name").
				SetDescription("Set ServiceBus shared access key name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key").
				SetTitle("Access Key").
				SetDescription("Set ServiceBus shared access key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("queue_name").
				SetTitle("Queue").
				SetDescription("Set ServiceBus queue name").
				SetMust(true).
				SetDefault(""),
		)
}
