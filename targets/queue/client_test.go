package queue

import (
	"context"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"

	"github.com/stretchr/testify/require"
	"testing"

	"time"
)

type mockQueueReceiver struct {
	host    string
	port    int
	channel string
	timeout int32
}

func (m *mockQueueReceiver) run(ctx context.Context) (*kubemq.QueueMessage, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(m.host, m.port),
		kubemq.WithClientId("response-id"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		return nil, err
	}

	queueMessages, err := client.ReceiveQueueMessages(ctx, &kubemq.ReceiveQueueMessagesRequest{
		RequestID:           "id",
		ClientID:            nuid.Next(),
		Channel:             m.channel,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     m.timeout,
		IsPeak:              false,
	})
	if err != nil {
		return nil, err
	}
	if len(queueMessages.Messages) == 0 {
		return nil, nil
	}
	return queueMessages.Messages[0], nil
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.Spec
		mockReceiver *mockQueueReceiver
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
					"channel": "queues",
				},
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues",
				timeout: 5,
			},
			sendReq: types.NewRequest().
				SetData([]byte("data")),
			wantReq: types.NewRequest().
				SetData([]byte("data")),
			wantResp: types.NewResponse(),

			wantErr: false,
		},
		{
			name: "bad request - bad channel",
			cfg: config.Spec{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
					"channel": "queues>  ",
				},
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues",
				timeout: 5,
			},
			sendReq: types.NewRequest().
				SetData([]byte("data")),
			wantReq: types.NewRequest().
				SetData([]byte("data")),
			wantResp: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			recRequestCh := make(chan *kubemq.QueueMessage, 1)
			recErrCh := make(chan error, 1)
			go func() {
				gotRequest, err := tt.mockReceiver.run(ctx)
				select {
				case recErrCh <- err:
				case recRequestCh <- gotRequest:
				}
			}()
			time.Sleep(time.Second)
			target := New()
			err := target.Init(ctx, tt.cfg)
			require.NoError(t, err)
			gotResp, err := target.Do(ctx, tt.sendReq)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
			select {
			case gotRequest := <-recRequestCh:
				require.EqualValues(t, tt.wantReq.Data, gotRequest.Body)
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
					"host": "localhost",
					"port": "-1",
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
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
