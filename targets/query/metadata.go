package query

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/types"
	"math"
	"time"
)

type metadata struct {
	id       string
	metadata string
	channel  string
	timeout  time.Duration
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	m.id = meta.ParseString("id", "")
	m.metadata = meta.ParseString("metadata", "")
	m.channel = meta.ParseString("channel", opts.defaultChannel)
	if m.channel == "" {
		return metadata{}, fmt.Errorf("channel cannot be empty")
	}
	timout, err := meta.ParseIntWithRange("timeout_seconds", opts.defaultTimeoutSeconds, 1, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing timeout seconds, %w", err)
	}
	m.timeout = time.Duration(timout) * time.Second
	return m, nil
}
