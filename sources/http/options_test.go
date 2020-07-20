package http

import (
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptions_parseOptions(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "valid options",
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "8080",
					"path": "",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad port",
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "-1",
					"path": "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseOptions(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}
