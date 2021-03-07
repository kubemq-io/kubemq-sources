package mqtt

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("messaging.mqtt").
		SetDescription("MQTT source properties").
		SetName("MQTT").
		SetProvider("").
		SetCategory("Messaging").
		SetTags("pub/sub", "iot").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetDescription("Set MQTT Host:Port connection").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("topic").
				SetDescription("Set subscribed topic name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("dynamic_mapping").
				SetDescription("Set Topic/Channel dynamic mapping").
				SetMust(true).
				SetDefault("true"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("client_id").
				SetTitle("Client ID").
				SetDescription("Set client ID").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("qos").
				SetTitle("QOS Level").
				SetDescription("Set QoS level").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(2),
		)
}
