package activemq

import (
	"context"
	"fmt"
	"github.com/go-stomp/stomp"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/targets"
	"github.com/kubemq-hub/kubemq-sources/targets/null"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func setupClient(ctx context.Context, queue string, target middleware.Middleware) (*Client, error) {
	c := New()
	err := c.Init(ctx, config.Metadata{
		Name: "activemq",
		Kind: "",
		Properties: map[string]string{
			"host":        "localhost:61613",
			"destination": queue,
			"username":    "admin",
			"password":    "admin",
		},
	})
	if err != nil {
		return nil, err
	}
	err = c.Start(ctx, target)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	return c, nil
}

func sendMessage(queue string, data []byte) error {

	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login("admin", "admin"),
		stomp.ConnOpt.Host("/"),
	}
	conn, err := stomp.Dial("tcp", "localhost:61613", options...)
	if err != nil {
		return fmt.Errorf("error connecting to activemq broker, %w", err)
	}
	return conn.Send(queue, "text/plain", data)
}

func TestClient_Start(t *testing.T) {
	tests := []struct {
		name    string
		target  targets.Target
		req     *types.Request
		queue   string
		wantErr bool
	}{
		{
			name: "request",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:     types.NewRequest().SetData([]byte("some-data")),
			queue:   "some-queue",
			wantErr: false,
		},
		{
			name: "request with target error",
			target: &null.Client{
				Delay:         0,
				DoError:       fmt.Errorf("some-error"),
				ResponseError: nil,
			},
			req:     types.NewRequest().SetData([]byte("some-data")),
			queue:   "some-queue",
			wantErr: false,
		},
		{
			name:    "request with nil target",
			target:  nil,
			req:     types.NewRequest().SetData([]byte("some-data")),
			queue:   "some-queue",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c, err := setupClient(ctx, tt.queue, tt.target)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			defer func() {
				_ = c.Stop()
			}()

			err = sendMessage(tt.queue, tt.req.Data)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
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
				Name: "activemq",
				Kind: "",
				Properties: map[string]string{
					"host":        "localhost:61613",
					"destination": "some-queue",
					"username":    "admin",
					"password":    "admin",
				},
			},
			wantErr: false,
		},
		{
			name: "init - bad url",
			cfg: config.Metadata{
				Name: "activemq-target",
				Kind: "",
				Properties: map[string]string{
					"host":        "localhost:2000",
					"destination": "some-queue",
					"username":    "admin",
					"password":    "admin",
				},
			},
			wantErr: true,
		},
		{
			name: "bad init - no  url",
			cfg: config.Metadata{
				Name: "activemq",
				Kind: "",
				Properties: map[string]string{
					"destination": "some-queue",
					"username":    "admin",
					"password":    "admin",
				},
			},
			wantErr: true,
		},
		{
			name: "bad init - no destination",
			cfg: config.Metadata{
				Name: "activemq",
				Kind: "",
				Properties: map[string]string{
					"host":     "localhost:61613",
					"username": "admin",
					"password": "admin",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}
