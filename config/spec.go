package config

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/types"
)

type Spec struct {
	Name       string            `json:"name"`
	Kind       string            `json:"kind"`
	Properties types.Metadata `json:"properties"`
}

func (s Spec) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if s.Kind == "" {
		return fmt.Errorf("kind cannot be empty")
	}
	return nil
}

