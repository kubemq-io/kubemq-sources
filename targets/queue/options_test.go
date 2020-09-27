package queue

import (
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/stretchr/testify/require"
	"testing"
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
				host:              "localhost",
				port:              50000,
				clientId:          "client_id",
				authToken:         "some-auth token",
				channel:           "some-channel",
				expirationSeconds: 0,
				delaySeconds:      0,
				maxReceiveCount:   0,
				deadLetterQueue:   "",
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad port",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":    "localhost:-1",
					"client_id":  "client_id",
					"auth_token": "some-auth token",
					"channel":    "some-channel",
				},
			},
			wantOpts: options{},
			wantErr:  true,
		},
		{
			name: "invalid options - no channel",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":    "localhost:50000",
					"client_id":  "client_id",
					"auth_token": "some-auth token",
					"channel":    "",
				},
			},
			wantOpts: options{},
			wantErr:  true,
		},
		{
			name: "invalid options - bad expiration value",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":            "localhost:50000",
					"client_id":          "client_id",
					"auth_token":         "some-auth token",
					"channel":            "some-channel",
					"expiration_seconds": "-1",
				},
			},
			wantOpts: options{},
			wantErr:  true,
		},
		{
			name: "invalid options - bad delay value",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":       "localhost:50000",
					"client_id":     "client_id",
					"auth_token":    "some-auth token",
					"channel":       "some-channel",
					"delay_seconds": "-1",
				},
			},
			wantOpts: options{},
			wantErr:  true,
		},
		{
			name: "invalid options - bad max receive value",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":           "localhost:50000",
					"client_id":         "client_id",
					"auth_token":        "some-auth token",
					"channel":           "some-channel",
					"max_receive_count": "-1",
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
