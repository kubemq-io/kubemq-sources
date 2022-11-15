package events_store

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
					"address":    "localhost:50000",
					"client_id":  "client_id",
					"auth_token": "some-auth token",
					"channel":    "some-channel",
				},
			},
			wantOpts: options{
				host:      "localhost",
				port:      50000,
				clientId:  "client_id",
				authToken: "some-auth token",
				channel:   "some-channel",
			},
			wantErr: false,
		},
		{
			name: "invalid options - no default channel",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":    "localhost:50000",
					"client_id":  "client_id",
					"auth_token": "some-auth token",
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
