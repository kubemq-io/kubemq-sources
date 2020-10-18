package kinesis

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("source.aws.kinesis").
		SetDescription("AWS Kinesis source properties").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("awsKey").
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
				SetDescription("Set Region").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("consumer_arn").
				SetDescription("Set Customer ARN").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shard_iterator_type").
				SetDescription("Set Shard Iterator Type").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shard_id").
				SetDescription("Set Shard Id").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("token").
				SetDescription("Set Token").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("sequence_number").
				SetDescription("Set Sequence Number").
				SetMust(false).
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
		)

}
