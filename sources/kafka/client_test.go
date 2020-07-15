package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/kubemq-hub/kubemq-source-connectors/middleware"
	"github.com/kubemq-hub/kubemq-source-connectors/types"
	"github.com/stretchr/testify/require"
)

type MockMiddleware struct {
}

func (m *MockMiddleware) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	return &types.Response{}, nil
}

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "kafka-target",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topics":  "TestTopicA,TestTopicB",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantExecErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		target  middleware.Middleware
		req     *types.Request
		wantErr bool
	}{
		{
			name: "valid publish request ",
			cfg: config.Metadata{
				Name: "kafka-target",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topics":  "TestTopic",
				},
			},

			req:     types.NewRequest().SetData([]byte("some-data")),
			target:  &MockMiddleware{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			err = c.Start(ctx, tt.target)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			for {
				if ctx.Err() != nil {
					return
				}

			}
			require.NoError(t, err)

		})
	}
}
