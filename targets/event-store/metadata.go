package event_store

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/types"
)

type metadata struct {
	id       string
	metadata string
	channel  string
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	m.id = meta.ParseString("id", "")
	m.metadata = meta.ParseString("metadata", "")
	m.channel = meta.ParseString("channel", opts.defaultChannel)
	if m.channel == "" {
		return metadata{}, fmt.Errorf("channel cannot be empty")
	}

	return m, nil
}
