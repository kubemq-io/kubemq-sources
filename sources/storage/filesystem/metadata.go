package filesystem

import (
	"encoding/json"
)

type TargetsMetadata map[string]string

func NewTargetMetadata() TargetsMetadata {
	return map[string]string{}
}

func (m TargetsMetadata) String() string {
	str, err := json.Marshal(&m)
	if err != nil {
		return ""
	}
	return string(str)
}

func (m TargetsMetadata) Set(key, value string) TargetsMetadata {
	m[key] = value
	return m
}
