package queue

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/types"
	"math"
)

type metadata struct {
	id                string
	metadata          string
	channel           string
	expirationSeconds int
	delaySeconds      int
	maxReceiveCount   int
	deadLetterQueue   string
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	m.id = meta.ParseString("id", "")
	m.metadata = meta.ParseString("metadata", "")
	m.channel = meta.ParseString("channel", opts.defaultChannel)
	if m.channel == "" {
		return metadata{}, fmt.Errorf("channel cannot be empty")
	}
	var err error
	m.expirationSeconds, err = meta.ParseIntWithRange("expiration_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing expiration seconds, %w", err)
	}
	m.delaySeconds, err = meta.ParseIntWithRange("delay_seconds", 0, 0, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing delay seconds, %w", err)
	}
	m.maxReceiveCount, err = meta.ParseIntWithRange("max_receive_count", 1, 1, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error max receive count seconds")
	}
	m.deadLetterQueue = meta.ParseString("dead_letter_queue", "")
	return m, nil
}
