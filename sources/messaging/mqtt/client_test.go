package mqtt

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
	m.channelName = "messaging.mqtt"
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
				Name: "messaging-mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"topic":     "some-queue",
					"username":  "",
					"password":  "",
					"client_id": nuid.Next(),
					"qos":       "0",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - bad url",
			cfg: config.Spec{
				Name: "messaging-mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:2000",
					"topic":     "some-queue",
					"username":  "",
					"password":  "",
					"client_id": nuid.Next(),
					"qos":       "0",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - no  url",
			cfg: config.Spec{
				Name: "messaging- mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"topic":     "some-queue",
					"username":  "",
					"password":  "",
					"client_id": nuid.Next(),
					"qos":       "0",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - no topic",
			cfg: config.Spec{
				Name: "messaging- mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"username":  "",
					"password":  "",
					"client_id": nuid.Next(),
					"qos":       "0",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - bad qos",
			cfg: config.Spec{
				Name: "messaging-mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"topic":     "some-queue",
					"username":  "",
					"password":  "",
					"client_id": nuid.Next(),
					"qos":       "-1",
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
			name: "valid mqtt receive",
			cfg: config.Spec{
				Name: "messaging-mqtt",
				Kind: "messaging.mqtt",
				Properties: map[string]string{
					"host":      "localhost:1883",
					"topic":     "some-queue",
					"username":  "",
					"password":  "",
					"client_id": nuid.Next(),
					"qos":       "0",
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
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			err = c.Start(ctx, tt.middleware)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			time.Sleep(time.Duration(30) * time.Second)
			require.NoError(t, err)
		})
	}
}
