package builder

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
	"strings"
)

type ConnectorConfig struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Type   string `json:"type"`
		Config string `json:"config"`
	} `json:"spec"`
}

func GetBuildManifest(url string) (*ConnectorConfig, error) {

	resp, err := resty.New().NewRequest().Get(url)
	if err != nil {
		return nil, err
	}
	manifests := string(resp.Body())
	for _, manifest := range strings.Split(manifests, "\n---\n") {
		c, err := parseConnectorConfig(manifest)
		if err != nil {
			continue
		}
		return c, nil
	}
	return nil, fmt.Errorf("no valid connector configuration found ")
}

func parseConnectorConfig(manifest string) (*ConnectorConfig, error) {
	c := &ConnectorConfig{}
	err := yaml.Unmarshal([]byte(manifest), c)
	if err != nil {
		return nil, err
	}
	switch c.Spec.Type {
	case "bridges", "targets", "sources":
		return c, nil
	default:
		return nil, fmt.Errorf("config is not a connector")
	}

}
