package activemq

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
	m.channelName = "event.messaging.activemq"
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
				Name: "activemq",
				Kind: "",
				Properties: map[string]string{
					"host":            "localhost:61613",
					"destination":     "test",
					"username":        "admin",
					"password":        "admin",
					"dynamic_mapping": "false",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - bad url",
			cfg: config.Spec{
				Name: "activemq-target",
				Kind: "",
				Properties: map[string]string{
					"host":            "localhost:8161",
					"destination":     "some-queue",
					"username":        "admin",
					"password":        "admin",
					"dynamic_mapping": "false",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - no  url",
			cfg: config.Spec{
				Name: "activemq",
				Kind: "",
				Properties: map[string]string{
					"destination":     "some-queue",
					"username":        "admin",
					"password":        "admin",
					"dynamic_mapping": "false",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - no destination",
			cfg: config.Spec{
				Name: "activemq",
				Kind: "",
				Properties: map[string]string{
					"host":            "localhost:61613",
					"username":        "admin",
					"password":        "admin",
					"dynamic_mapping": "false",
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
			name: "valid activemq receive",
			cfg: config.Spec{
				Name: "messaging-activemq",
				Kind: "messaging.activemq",
				Properties: map[string]string{
					"host":            "localhost:61613",
					"destination":     "some-queue",
					"username":        "admin",
					"password":        "admin",
					"dynamic_mapping": "false",
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
			time.Sleep(time.Duration(30) * time.Second)
			require.NoError(t, err)
		})
	}
}
