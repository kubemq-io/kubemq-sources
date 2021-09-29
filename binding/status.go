package binding

import (
	"fmt"
	"github.com/kubemq-io/kubemq-sources/config"
)

type Status struct {
	Binding          string            `json:"binding"`
	Ready            bool              `json:"ready"`
	SourceType       string            `json:"source_type"`
	SourceConfig     map[string]string `json:"source_config"`
	TargetType       string            `json:"target_type"`
	TargetConnection string            `json:"target_connection"`
	TargetConfig     map[string]string `json:"target_config"`
}

func getTargetConnection(properties map[string]string) string {
	return fmt.Sprintf("%s/%s", properties["address"], properties["channel"])
}
func newStatus(cfg config.BindingConfig) *Status {
	return &Status{
		Binding:          cfg.Name,
		Ready:            false,
		SourceType:       cfg.Source.Kind,
		SourceConfig:     cfg.Source.Properties,
		TargetType:       cfg.Target.Kind,
		TargetConnection: getTargetConnection(cfg.Target.Properties),
		TargetConfig:     cfg.Target.Properties,
	}
}
