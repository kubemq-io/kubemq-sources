package query

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/types"

	"github.com/stretchr/testify/require"
	"testing"

	"time"
)

type mockQueryReceiver struct {
	host           string
	port           int
	channel        string
	executionDelay time.Duration
	executionError error
	executionTime  int64
}

func (m *mockQueryReceiver) run(ctx context.Context, t *testing.T) error {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(m.host, m.port),
		kubemq.WithClientId("response-id"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		return err
	}
	errCh := make(chan error, 1)
	queryCh, err := client.SubscribeToQueries(ctx, m.channel, "", errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case query := <-queryCh:
				time.Sleep(m.executionDelay)
				queryResponse := client.R().SetRequestId(query.Id).SetResponseTo(query.ResponseTo).SetExecutedAt(time.Unix(m.executionTime, 0))
				if m.executionError != nil {
					queryResponse.SetError(m.executionError)
				} else {
					queryResponse.SetBody(query.Body)
					queryResponse.SetMetadata(query.Metadata)
				}
				err := queryResponse.Send(ctx)
				require.NoError(t, err)
			case err := <-errCh:
				require.NoError(t, err)
			case <-ctx.Done():
				return
			}

		}
	}()
	time.Sleep(time.Second)
	return nil
}
func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.Spec
		mockReceiver *mockQueryReceiver
		req          *types.Request
		wantResp     *types.Response
		wantErr      bool
	}{
		{
			name: "request",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":         "localhost:50000",
					"channel":         "queries",
					"timeout_seconds": "5",
				},
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req: types.NewRequest().
				SetData([]byte("data")).
				SetMetadata("id"),
			wantResp: types.NewResponse().
				SetMetadata("id").
				SetData([]byte("data")),
			wantErr: false,
		},
		{
			name: "request with execution error",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":         "localhost:50000",
					"channel":         "queries",
					"timeout_seconds": "5",
				},
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries",
				executionDelay: 0,
				executionError: fmt.Errorf("error"),
				executionTime:  0,
			},
			req: types.NewRequest().
				SetData([]byte("data")),
			wantResp: types.NewResponse().
				SetError(fmt.Errorf("error")),
			wantErr: false,
		},
		{
			name: "request error - timeout seconds",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":         "localhost:50000",
					"channel":         "queries",
					"timeout_seconds": "1",
				},
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries",
				executionDelay: 3 * time.Second,
				executionError: nil,
				executionTime:  0,
			},
			req: types.NewRequest().
				SetMetadata(`"id": "id"`),
			wantResp: types.NewResponse().SetError(fmt.Errorf("rpc error: code = Internal desc = Error 301: timeout for request message")),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := tt.mockReceiver.run(ctx, t)
			require.NoError(t, err)
			target := New()
			err = target.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			gotResp, err := target.Do(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
		})
	}
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
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
			wantErr: false,
		},
		{
			name: "init - error",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:-1",
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

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
