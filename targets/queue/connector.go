package queue

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.queue").
		SetDescription("Kubemq Queue Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetDescription("Set Kubemq grpc endpoint address").
				SetMust(true).
				SetDefault("").
				SetLoadedOptions("kubemq-address"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("connection-type").
				SetDescription("Set Channel Mapping mode").
				SetMust(true).
				SetOptions([]string{"Implicit", "Dynamic"}).
				SetDefault("Implicit").
				NewCondition("Implicit", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("dynamic_mapping").
						SetDescription("Set dynamic mapping").
						SetMust(true).
						SetDefault("false"),
					common.NewProperty().
						SetKind("string").
						SetName("channel").
						SetDescription("Set Events channel").
						SetMust(true).
						SetDefaultFromKey("channel.queue"),
				}).
				NewCondition("Dynamic", []*common.Property{
					common.NewProperty().
						SetKind("null").
						SetName("dynamic_mapping").
						SetDescription("Set dynamic mapping").
						SetMust(true).
						SetDefault("true"),
				}),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Set Queue connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetDescription("Set Queue connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("expiration_seconds").
				SetDescription("Set Queue message expiration in seconds").
				SetMust(false).
				SetMin(0).
				SetMax(math.MaxInt32).
				SetDefault("0"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("delay_seconds").
				SetDescription("Set Queue message delay in seconds").
				SetMust(false).
				SetMin(0).
				SetMax(math.MaxInt32).
				SetDefault("0"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("dead_letter_queue").
				SetDescription("Set Queue message dead-letter queue name").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_receive_count").
				SetDescription("Set Queue message max fails retries route to dead-letter").
				SetMust(false).
				SetMin(0).
				SetMax(math.MaxInt32).
				SetDefault("0"),
		)
}
