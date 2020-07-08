package targets

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
			name: "http target",
			cfgs: []config.Metadata{
				{
					Name: "http-1",
					Kind: "target.http",
					Properties: map[string]string{
						"uri":       "",
						"auth_type": "no_auth",
						"username":  "",
						"password":  "",
						"token":     "",
						"headers":   ``,
					},
				},
			},
			want: map[string]string{
				"http-1": "http-1",
			},
			wantErr: false,
		},
		{
			name: "http target - bad init",
			cfgs: []config.Metadata{
				{
					Name: "http-1",
					Kind: "target.http",
					Properties: map[string]string{
						"headers": `bad map`,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "http target - duplicate names",
			cfgs: []config.Metadata{
				{
					Name: "http-1",
					Kind: "target.http",
					Properties: map[string]string{
						"uri":       "",
						"auth_type": "no_auth",
						"username":  "",
						"password":  "",
						"token":     "",
						"headers":   ``,
					},
				},
				{
					Name: "http-1",
					Kind: "target.http",
					Properties: map[string]string{
						"uri":       "",
						"auth_type": "no_auth",
						"username":  "",
						"password":  "",
						"token":     "",
						"headers":   ``,
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
					Name:       "target-1",
					Kind:       "some-bad-target",
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
			targets := make(map[string]string)
			for key, source := range got {
				targets[key] = source.Name()
			}
			if !reflect.DeepEqual(targets, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
