package events_store

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.events-store").
		SetDescription("Kubemq Events-Store Target").
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
				SetDescription("Set Events-Store channel").
				SetMust(true).
				SetDefault("events-store"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Set Events-Store connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetDescription("Set Events-Store connection authentication token").
				SetMust(false).
				SetDefault(""),
		)
}
