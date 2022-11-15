package rabbitmq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/types"
	"github.com/nats-io/nuid"
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
	m.channelName = "events.messaging.rabbitmq"
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
			name: "init",
			cfg: config.Spec{
				Name: "messaging-rabbitmq",
				Kind: "messaging.rabbitmq",
				Properties: map[string]string{
					"url":              "amqp://guest:guest@localhost:5672/",
					"queue":            "some-queue",
					"consumer":         nuid.Next(),
					"requeue_on_error": "false",
					"dynamic_mapping":  "false",
					"auto_ack":         "false",
					"exclusive":        "false",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - bad url",
			cfg: config.Spec{
				Name: "messaging-rabbitmq",
				Kind: "messaging.rabbitmq",
				Properties: map[string]string{
					"url":              "amqp://guest:guest@localhost:6000/",
					"consumer":         nuid.Next(),
					"queue":            "some-queue",
					"dynamic_mapping":  "false",
					"requeue_on_error": "false",
					"auto_ack":         "false",
					"exclusive":        "false",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid bad init - missing url",
			cfg: config.Spec{
				Name: "messaging-rabbitmq",
				Kind: "messaging.rabbitmq",
				Properties: map[string]string{
					"queue":            "some-queue",
					"consumer":         nuid.Next(),
					"requeue_on_error": "false",
					"dynamic_mapping":  "false",
					"auto_ack":         "false",
					"exclusive":        "false",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - missing queue",
			cfg: config.Spec{
				Name: "messaging-rabbitmq",
				Kind: "messaging.rabbitmq",
				Properties: map[string]string{
					"url":              "amqp://guest:guest@localhost:5672/",
					"consumer":         nuid.Next(),
					"requeue_on_error": "false",
					"auto_ack":         "false",
					"exclusive":        "false",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - missing consumer",
			cfg: config.Spec{
				Name: "messaging-rabbitmq",
				Kind: "messaging.rabbitmq",
				Properties: map[string]string{
					"url":              "amqp://guest:guest@localhost:5432/",
					"queue":            "some-queue",
					"requeue_on_error": "false",
					"auto_ack":         "false",
					"exclusive":        "false",
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
			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
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
		wantErr    bool
		middleware middleware.Middleware
	}{
		{
			name: "valid rabbitmq receive",
			cfg: config.Spec{
				Name: "messaging-rabbitmq",
				Kind: "messaging.rabbitmq",
				Properties: map[string]string{
					"url":              "amqp://guest:guest@localhost:5672/",
					"queue":            "some-queue",
					"consumer":         nuid.Next(),
					"dynamic_mapping":  "false",
					"requeue_on_error": "false",
					"auto_ack":         "false",
					"exclusive":        "false",
				},
			},
			middleware: middle,

			wantErr: false,
		}, {
			name: "invalid valid rabbitmq receive - fake queue",
			cfg: config.Spec{
				Name: "messaging-rabbitmq",
				Kind: "messaging.rabbitmq",
				Properties: map[string]string{
					"url":              "amqp://guest:guest@localhost:5672/",
					"queue":            "fake-queue",
					"consumer":         nuid.Next(),
					"dynamic_mapping":  "false",
					"requeue_on_error": "false",
					"auto_ack":         "false",
					"exclusive":        "false",
				},
			},
			middleware: middle,

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			err = c.Start(ctx, tt.middleware)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.Nil(t, err)
			time.Sleep(time.Duration(5) * time.Second)
			err = c.Stop()
			require.Nil(t, err)
			time.Sleep(time.Duration(15) * time.Second)
		})
	}
}
