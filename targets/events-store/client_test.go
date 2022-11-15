package events_store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/types"

	"github.com/stretchr/testify/require"
)

type mockEventStoreReceiver struct {
	host    string
	port    int
	channel string
	timeout time.Duration
}

func (m *mockEventStoreReceiver) run(ctx context.Context) (*kubemq.EventStoreReceive, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(m.host, m.port),
		kubemq.WithClientId("response-id"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		return nil, err
	}
	errCh := make(chan error, 1)
	eventCh, err := client.SubscribeToEventsStore(ctx, m.channel, "", errCh, kubemq.StartFromNewEvents())
	if err != nil {
		return nil, err
	}
	select {
	case event := <-eventCh:
		return event, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, nil
	case <-time.After(m.timeout):
		return nil, fmt.Errorf("timeout")
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.Spec
		mockReceiver *mockEventStoreReceiver
		sendReq      *types.Request
		wantReq      *types.Request
		wantResp     *types.Response
		wantErr      bool
	}{
		{
			name: "request",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
					"channel": "event_stores",
				},
			},
			mockReceiver: &mockEventStoreReceiver{
				host:    "localhost",
				port:    50000,
				channel: "event_stores",
				timeout: 5 * time.Second,
			},
			sendReq: types.NewRequest().
				SetData([]byte("data")),
			wantReq: types.NewRequest().
				SetData([]byte("data")),
			wantResp: types.NewResponse(),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			recEventCh := make(chan *kubemq.EventStoreReceive, 1)
			recErrCh := make(chan error, 1)
			go func() {
				gotEvent, err := tt.mockReceiver.run(ctx)
				select {
				case recErrCh <- err:
				case recEventCh <- gotEvent:
				}
			}()
			time.Sleep(time.Second)
			target := New()
			err := target.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			gotResp, err := target.Do(ctx, tt.sendReq)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
			select {
			case gotEvent := <-recEventCh:
				require.EqualValues(t, tt.wantReq.Data, gotEvent.Body)
			case err := <-recErrCh:
				require.NoError(t, err)
			case <-ctx.Done():
				require.NoError(t, ctx.Err())
			}
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
					"address":    "localhost:50000",
					"client_id":  "client_id",
					"auth_token": "some-auth token",
					"channel":    "some-channel",
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
