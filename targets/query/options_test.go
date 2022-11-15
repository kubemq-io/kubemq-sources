package query

import (
	"testing"

	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/stretchr/testify/require"
)

func TestOptions_parseOptions(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.Spec
		wantOpts options
		wantErr  bool
	}{
		{
			name: "valid options",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":         "localhost:50000",
					"client_id":       "client_id",
					"auth_token":      "some-auth token",
					"channel":         "some-channel",
					"timeout_seconds": "100",
				},
			},
			wantOpts: options{
				host:           "localhost",
				port:           50000,
				clientId:       "client_id",
				authToken:      "some-auth token",
				channel:        "some-channel",
				timeoutSeconds: 100,
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad port",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:-1",
				},
			},
			wantOpts: options{},
			wantErr:  true,
		},
		{
			name: "invalid options - no default channel",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":         "localhost:50000",
					"client_id":       "client_id",
					"auth_token":      "some-auth token",
					"timeout_seconds": "100",
				},
			},
			wantOpts: options{},
			wantErr:  true,
		},
		{
			name: "invalid options - bad default timeout seconds",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":         "localhost:50000",
					"client_id":       "client_id",
					"auth_token":      "some-auth token",
					"channel":         "some-channel",
					"timeout_seconds": "-1",
				},
			},
			wantOpts: options{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOpts, err := parseOptions(tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, tt.wantOpts, gotOpts)
		})
	}
}
