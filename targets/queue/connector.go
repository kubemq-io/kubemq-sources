package queue

import (
	"github.com/kubemq-hub/builder/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.queue").
		SetDescription("Kubemq Queue Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("address").
				SetDescription("Sets Kubemq grpc endpoint address").
				SetMust(true).
				SetDefault("").
				SetLoadedOptions("kubemq-address"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("channel").
				SetDescription("Sets Queue channel").
				SetMust(true).
				SetDefault("queues"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Sets Queue connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetDescription("Sets Queue connection authentication token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("expiration_seconds").
				SetDescription("Sets Queue message expiration in seconds").
				SetMust(false).
				SetMin(0).
				SetMax(math.MaxInt32).
				SetDefault("0"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("delay_seconds").
				SetDescription("Sets Queue message delay in seconds").
				SetMust(false).
				SetMin(0).
				SetMax(math.MaxInt32).
				SetDefault("0"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("dead_letter_queue").
				SetDescription("Sets Queue message dead-letter queue name").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_receive_count").
				SetDescription("Sets Queue message max fails retries route to dead-letter").
				SetMust(false).
				SetMin(0).
				SetMax(math.MaxInt32).
				SetDefault("0"),
		)
}
