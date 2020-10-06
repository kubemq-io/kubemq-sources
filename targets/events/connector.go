package events

import "github.com/kubemq-hub/builder/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("target.events").
		SetDescription("Kubemq Events Target").
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
				SetDescription("Sets Events channel").
				SetMust(true).
				SetDefault("events"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetDescription("Sets Events connection client Id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("auth_token").
				SetDescription("Sets Events connection authentication token").
				SetMust(false).
				SetDefault(""),
		)
}
