package sqs

import (
	"math"

	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.sqs").
		SetDescription("AWS SQS source properties").
		SetName("SQS").
		SetProvider("AWS").
		SetCategory("Messaging").
		SetTags("queue", "cloud", "managed").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set AWS Key").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set AWS Secret Key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set AWS Region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set AWS Token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("queue").
				SetDescription("Set Sqs Queue name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("pull_delay").
				SetDescription("Set Pull Delay in seconds").
				SetMust(false).
				SetDefault("5").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("visibility_timeout").
				SetDescription("Set Visibility Timout").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_number_of_messages").
				SetTitle("Max Messages").
				SetDescription("Set Max Number of Messages").
				SetMust(false).
				SetDefault("1").
				SetMin(0).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("wait_time_seconds").
				SetTitle("Wait time (seconds)").
				SetDescription("Set Wait Time Second").
				SetMust(false).
				SetDefault("0").
				SetMin(0).
				SetMax(math.MaxInt32),
		)
}
