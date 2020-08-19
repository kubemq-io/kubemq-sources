package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/stretchr/testify/require"
)

type MockMiddleware struct {
	error     chan error
	wantError bool
}

func (m *MockMiddleware) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	if m.wantError {
		err := fmt.Errorf("newError")
		m.error <- err
		return nil, err
	}

	m.error <- nil
	return &types.Response{
		Data: request.Data,
	}, nil

}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "valid init",
			cfg: config.Metadata{
				Name: "kafka-target",
				Properties: map[string]string{
					"brokers":       "localhost:9092",
					"topics":        "TestTopicA,TestTopicB",
					"consumerGroup": "test_client",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init",
			cfg: config.Metadata{
				Name: "kafka-target",
				Properties: map[string]string{
					"brokers":       "localhost:9090",
					"topics":        "TestTopic",
					"consumerGroup": "test_client1",
				},
			},
			wantErr: true,
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
	errors := make(chan error)
	tests := []struct {
		name         string
		cfg          config.Metadata
		target       middleware.Middleware
		req          *types.Request
		wantErr      bool
		wantErrorMsg bool
	}{
		{
			name: "valid connection target error ",
			cfg: config.Metadata{
				Name: "kafka-target",
				Properties: map[string]string{
					"brokers":       "localhost:9092",
					"topics":        "TestTopic",
					"consumerGroup": "test_client1",
				},
			},

			req: types.NewRequest().SetData([]byte("some-data")),
			target: &MockMiddleware{
				error:     errors,
				wantError: true,
			},
			wantErr:      false,
			wantErrorMsg: true,
		},
		{
			name: "valid connection target success ",
			cfg: config.Metadata{
				Name: "kafka-target",
				Properties: map[string]string{
					"brokers":       "localhost:9092",
					"topics":        "TestTopic",
					"consumerGroup": "test_client1",
				},
			},
			req:          types.NewRequest().SetData([]byte("some-data")),
			wantErrorMsg: false,
			target: &MockMiddleware{
				error:     errors,
				wantError: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			end := make(chan bool)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			err = c.Start(ctx, tt.target)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			go func() {
				for {
					if ctx.Err() != nil {
						end <- true
					}
				}
			}()
			err = <-errors
			if tt.wantErrorMsg {
				require.Error(t, err)
				return
			}
			require.Nil(t, err)
			return
		})

	}
}
