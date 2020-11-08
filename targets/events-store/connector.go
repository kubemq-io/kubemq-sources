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
						SetDefaultFromKey("channel.events-store"),
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
