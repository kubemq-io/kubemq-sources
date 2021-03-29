package command

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.command").
		SetDescription("Kubemq Command Target").
		SetName("KubeMQ Command").
		SetProvider("").
		SetCategory("RPC").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetTitle("KubeMQ gRPC Service Address").
				SetDescription("Set Kubemq grpc endpoint address").
				SetMust(true).
				SetDefault("kubemq-cluster-grpc.kubemq:50000").
				SetLoadedOptions("kubemq-address"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("channel").
				SetDescription("Set Events channel").
				SetMust(true).
				SetDefaultFromKey("channel.command"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetTitle("Client ID").
				SetDescription("Set Command connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetTitle("Authentication Token").
				SetDescription("Set Command connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("timeout_seconds").
				SetTitle("Timeout (Seconds)").
				SetDescription("Set Command request timeout in seconds").
				SetMust(false).
				SetDefault("60").
				SetMin(1).
				SetMax(math.MaxInt32),
		)
}
