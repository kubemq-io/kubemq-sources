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
				SetKind("string").
				SetName("channel").
				SetDescription("Set Command channel").
				SetMust(true).
				SetDefault("commands"),
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
