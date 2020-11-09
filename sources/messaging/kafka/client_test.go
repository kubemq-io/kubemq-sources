package kafka

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"testing"
	"time"

	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/stretchr/testify/require"
)

type mockMiddleware struct {
	client      *kubemq.Client
	channelName string
}

func (m *mockMiddleware) Init() {

	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithCheckConnection(true),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))

	if err != nil {
		panic(err)
	}
	m.client = client
	m.channelName = "event.messaging.kafka"
}

func (m *mockMiddleware) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	fmt.Println(request)
	r := types.NewResponse()
	r.SetData([]byte("ok"))
	r.SetMetadata(`"result":"ok"`)
	event := m.client.NewEvent()
	event.Channel = m.channelName
	event.Body = request.Data
	err := event.Send(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "valid init",
			cfg: config.Spec{
				Name: "messaging.kafka",
				Properties: map[string]string{
					"brokers":        "localhost:9092",
					"topics":         "TestTopicA,TestTopicB",
					"dynamic_mapping":  "false",
					"consumer_group": "test_client",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - incorrect brokers ",
			cfg: config.Spec{
				Name: "messaging.kafka",
				Properties: map[string]string{
					"brokers":        "localhost:9090",
					"topics":         "TestTopic",
					"dynamic_mapping":  "false",
					"consumer_group": "test_client1",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing brokers ",
			cfg: config.Spec{
				Name: "messaging.kafka",
				Properties: map[string]string{
					"topics":         "TestTopic",
					"consumer_group": "test_client1",
					"dynamic_mapping":  "false",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing topics ",
			cfg: config.Spec{
				Name: "messaging.kafka",
				Properties: map[string]string{
					"brokers":        "localhost:9090",
					"consumer_group": "test_client1",
					"dynamic_mapping":  "false",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing consumer_group ",
			cfg: config.Spec{
				Name: "messaging.kafka",
				Properties: map[string]string{
					"brokers": "localhost:9092",
					"topics":  "TestTopic",
					"dynamic_mapping":  "false",
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

		})
	}
}

func TestClient_Do(t *testing.T) {
	middle := &mockMiddleware{}
	middle.Init()
	tests := []struct {
		name       string
		cfg        config.Spec
		middleware middleware.Middleware
		req        *types.Request
		wantErr    bool
	}{
		{
			name: "valid connection target",
			cfg: config.Spec{
				Name: "messaging.kafka",
				Properties: map[string]string{
					"brokers":        "localhost:9092",
					"topics":         "TestTopic",
					"consumer_group": "test_client1",
					"dynamic_mapping":  "false",
				},
			},

			req:        types.NewRequest().SetData([]byte("some-data")),
			middleware: middle,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			err = c.Start(ctx, tt.middleware)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.Nil(t, err)
			time.Sleep(time.Duration(1) * time.Second)
			err = c.Stop()
			require.Nil(t, err)
			time.Sleep(time.Duration(15) * time.Second)

		})

	}
}
