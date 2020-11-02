package eventhubs

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("azure.eventhubs").
		SetDescription("Azure EventHubs Source").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("end_point").
				SetDescription("Set EventHubs end point").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key_name").
				SetDescription("Set EventHubs shared access key name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("shared_access_key").
				SetDescription("Set EventHubs shared access key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("entity_path").
				SetDescription("Set EventHubs entity path").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("partition_id").
				SetDescription("Set EventHubs partition id").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(common.NewProperty().
			SetKind("string").
			SetName("receive_type").
			SetOptions([]string{"latest_offset", "from_timestamp", "with_consumer_group", "with_epoch", "with_prefetch_count", "with_starting_offset"}).
			SetDescription("Set partition receive_type").
			SetMust(true).
			SetDefault("latest_offset")).
		AddProperty(
			common.NewProperty().
				SetKind("condition").
				SetName("receive_type").
				SetOptions([]string{"latest_offset", "from_timestamp", "with_consumer_group", "with_epoch", "with_prefetch_count", "with_starting_offset"}).
				SetDescription("Set partition conditions").
				SetMust(true).
				SetDefault("latest_offset").
				NewCondition("from_timestamp", []*common.Property{
					common.NewProperty().
						SetKind("int").
						SetName("time_stamp").
						SetDescription("Set timestamp to collect events from (RFC3339)").
						SetMust(true).
						SetDefault("").
						SetMin(0).
						SetMax(math.MaxInt64),
				}).
				NewCondition("with_consumer_group", []*common.Property{
					common.NewProperty().
						SetKind("string").
						SetName("consumer_group").
						SetDescription("Set the Consumer Group to collect events from").
						SetMust(true).
						SetDefault(""),
				}).
				NewCondition("with_epoch", []*common.Property{
					common.NewProperty().
						SetKind("int").
						SetName("epoch").
						SetDescription("Set timestamp to collect events from (epoch)").
						SetMust(true).
						SetDefault("").
						SetMin(0).
						SetMax(math.MaxInt64),
				}).
				NewCondition("with_prefetch_count", []*common.Property{
					common.NewProperty().
						SetKind("int").
						SetName("prefetch_count").
						SetDescription("Set Prefetch count to collect events from").
						SetMust(true).
						SetDefault("").
						SetMin(0).
						SetMax(math.MaxInt64),
				}).
				NewCondition("with_starting_offset", []*common.Property{
					common.NewProperty().
						SetKind("string").
						SetName("starting_offset").
						SetDescription("Set starting offset").
						SetMust(true).
						SetDefault(""),
				}),
		)
}
