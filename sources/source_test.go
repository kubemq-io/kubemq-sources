package sources

import (
	"context"
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		cfgs    []config.Metadata
		want    map[string]string
		wantErr bool
	}{
		{
			name: "rpc source",
			cfgs: []config.Metadata{
				{
					Name: "rpc-1",
					Kind: "source.rpc",
					Properties: map[string]string{
						"host":                       "localhost",
						"port":                       "50000",
						"client_id":                  "",
						"auth_token":                 "some-auth token",
						"channel":                    "some-channel",
						"group":                      "",
						"concurrency":                "1",
						"auto_reconnect":             "true",
						"reconnect_interval_seconds": "1",
						"max_reconnects":             "0",
					},
				},
			},
			want: map[string]string{
				"rpc-1": "rpc-1",
			},
			wantErr: false,
		},
		{
			name: "rpc source - bad init",
			cfgs: []config.Metadata{
				{
					Name: "rpc-1",
					Kind: "source.rpc",
					Properties: map[string]string{
						"host": "localhost",
						"port": "-1",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "rpc source - duplicate names",
			cfgs: []config.Metadata{
				{
					Name: "rpc-1",
					Kind: "source.rpc",
					Properties: map[string]string{
						"host":                       "localhost",
						"port":                       "50000",
						"client_id":                  "",
						"auth_token":                 "some-auth token",
						"channel":                    "some-channel",
						"group":                      "",
						"concurrency":                "1",
						"auto_reconnect":             "true",
						"reconnect_interval_seconds": "1",
						"max_reconnects":             "0",
					},
				}, {
					Name: "rpc-1",
					Kind: "source.rpc",
					Properties: map[string]string{
						"host":                       "localhost",
						"port":                       "50000",
						"client_id":                  "",
						"auth_token":                 "some-auth token",
						"channel":                    "some-channel",
						"group":                      "",
						"concurrency":                "1",
						"auto_reconnect":             "true",
						"reconnect_interval_seconds": "1",
						"max_reconnects":             "0",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "queue source",
			cfgs: []config.Metadata{
				{
					Name: "queue-1",
					Kind: "source.queue",
					Properties: map[string]string{
						"host":             "localhost",
						"port":             "50000",
						"client_id":        "some-client-id",
						"auth_token":       "",
						"channel":          "queue",
						"response_channel": "default-response",
						"concurrency":      "1",
						"batch_size":       "1",
						"wait_timeout":     "60",
					},
				},
			},
			want: map[string]string{
				"queue-1": "queue-1",
			},
			wantErr: false,
		},
		{
			name: "queue source - bad init",
			cfgs: []config.Metadata{
				{
					Name: "queue-1",
					Kind: "source.queue",
					Properties: map[string]string{
						"host": "localhost",
						"port": "-1",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid source",
			cfgs: []config.Metadata{
				{
					Name:       "queue-1",
					Kind:       "some-bad-source",
					Properties: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			got, err := Load(ctx, tt.cfgs)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			sources := make(map[string]string)
			for key, source := range got {
				sources[key] = source.Name()
			}
			if !reflect.DeepEqual(sources, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
