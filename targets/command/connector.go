package command

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.command").
		SetDescription("Kubemq Command Target").
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
						SetDefaultFromKey("channel.command"),
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
				SetDescription("Set Command connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetDescription("Set Command connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("timeout_seconds").
				SetDescription("Set Command request timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
