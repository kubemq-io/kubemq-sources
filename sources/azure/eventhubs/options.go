package eventhubs

import (
	"fmt"
	"time"

	"github.com/Azure/azure-event-hubs-go/v3"
	"github.com/kubemq-io/kubemq-sources/config"
)

const (
	defaultReceiveType = ""
	defaultOffset      = ""
)

type options struct {
	connectionString string
	partitionID      string
	receiveType      eventhub.ReceiveOption
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.partitionID, err = cfg.Properties.MustParseString("partition_id")
	if err != nil {
		return options{}, fmt.Errorf("error parsing partition_id , %w", err)
	}
	endPoint, err := cfg.Properties.MustParseString("end_point")
	if err != nil {
		return options{}, fmt.Errorf("error parsing end_point , %w", err)
	}
	sharedAccessKeyName, err := cfg.Properties.MustParseString("shared_access_key_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shared_access_key_name , %w", err)
	}
	sharedAccessKey, err := cfg.Properties.MustParseString("shared_access_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing shared_access_key , %w", err)
	}
	entityPath, err := cfg.Properties.MustParseString("entity_path")
	if err != nil {
		return options{}, fmt.Errorf("error parsing entity_path , %w", err)
	}
	receiveType := cfg.Properties.ParseString("receive_type", defaultReceiveType)

	switch receiveType {
	case "latest_offset":
		o.receiveType = eventhub.ReceiveWithLatestOffset()
	case "from_timestamp":
		timeStamp, err := cfg.Properties.MustParseString("time_stamp")
		if err != nil {
			return options{}, fmt.Errorf("when using from_timestamp please set a time_stamp , %w", err)
		}
		t, err := time.Parse(time.RFC3339, timeStamp)
		if err != nil {
			return options{}, fmt.Errorf("failed to parse time using format  RFC3339 ,when using from_timestamp please set a time_stamp %w", err)
		}
		o.receiveType = eventhub.ReceiveFromTimestamp(t)
	case "with_consumer_group":
		group, err := cfg.Properties.MustParseString("consumer_group")
		if err != nil {
			return options{}, fmt.Errorf("when using with_consumer_group please set a consumer_group , %w", err)
		}
		o.receiveType = eventhub.ReceiveWithConsumerGroup(group)
	case "with_epoch":
		epoch, err := cfg.Properties.MustParseInt("epoch")
		if err != nil {
			return options{}, fmt.Errorf("when using with_epoch please set a epoch , %w", err)
		}
		o.receiveType = eventhub.ReceiveWithEpoch(int64(epoch))
	case "with_prefetch_count":
		prefetchCount, err := cfg.Properties.MustParseInt("prefetch_count")
		if err != nil {
			return options{}, fmt.Errorf("when using with_prefetch_count please set a prefetch_count , %w", err)
		}
		o.receiveType = eventhub.ReceiveWithPrefetchCount(uint32(prefetchCount))
	case "with_starting_offset":
		startingOffset := cfg.Properties.ParseString("starting_offset", defaultOffset)
		o.receiveType = eventhub.ReceiveWithStartingOffset(startingOffset)
	default:
		startingOffset := cfg.Properties.ParseString("starting_offset", defaultOffset)
		o.receiveType = eventhub.ReceiveWithStartingOffset(startingOffset)
	}

	o.connectionString = fmt.Sprintf("Endpoint=%s;SharedAccessKeyName=%s;SharedAccessKey=%s;EntityPath=%s", endPoint, sharedAccessKeyName, sharedAccessKey, entityPath)
	return o, nil
}
