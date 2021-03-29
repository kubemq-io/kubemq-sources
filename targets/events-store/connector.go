package events_store

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("kubemq.events-store").
		SetDescription("Kubemq Events-Store Target").
		SetName("KubeMQ Events Store").
		SetProvider("").
		SetCategory("Pub/Sub").
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
				SetDescription("Set Events Store channel").
				SetMust(true).
				SetDefaultFromKey("channel.events-store"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetTitle("Client ID").
				SetDescription("Set Events-Store connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetTitle("Authentication Token").
				SetDescription("Set Events-Store connection authentication token").
				SetMust(false).
				SetDefault(""),
		)
}
