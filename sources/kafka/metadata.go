package kafka

import (
	"encoding/json"
	"fmt"

	b64 "encoding/base64"

	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-hub/kubemq-source-connectors/types"
)

type metadata struct {
	Headers []kafka.RecordHeader
	Key     []byte
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	var err error
	err = m.parseHeaders(meta.ParseString("Headers", ""))
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing headers, %w", err)
	}
	k := meta.ParseString("Key", "")
	err = m.parseKey(k)
	if err != nil {
		return metadata{}, fmt.Errorf("error parsing Key, %w", err)
	}

	return m, nil
}

func (meta *metadata) parseHeaders(headers string) error {
	if len(headers) == 0 {
		return nil
	}
	var h []kafka.RecordHeader
	err := json.Unmarshal([]byte(headers), &h)
	if err != nil {
		return err
	}
	meta.Headers = h
	return nil
}

func (meta *metadata) parseKey(key string) error {
	if len(key) == 0 {
		return nil
	}
	var err error
	meta.Key, err = b64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	return nil
}
