package servicebus

import (
	"context"
	"fmt"
	"github.com/fortytw2/leaktest"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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
	m.channelName = "event.azure.servicebus"
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

type testStructure struct {
	endPoint            string
	sharedAccessKeyName string
	sharedAccessKey     string
	queueName           string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/azure/servicebus/endPoint.txt")
	if err != nil {
		return nil, err
	}
	t.endPoint = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/servicebus/sharedAccessKeyName.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKeyName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/servicebus/sharedAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKey = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/servicebus/queueName.txt")
	if err != nil {
		return nil, err
	}
	t.queueName = fmt.Sprintf("%s", dat)

	return t, nil
}

func TestClient_Init(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "azure-servicebus",
				Kind: "azure.servicebus",
				Properties: map[string]string{
					"queue_name":             dat.queueName,
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing queue_name",
			cfg: config.Spec{
				Name: "azure-servicebus",
				Kind: "azure.servicebus",
				Properties: map[string]string{
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing end_point",
			cfg: config.Spec{
				Name: "azure-servicebus",
				Kind: "azure.servicebus",
				Properties: map[string]string{
					"queue_name":             dat.queueName,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing shared_access_key_name",
			cfg: config.Spec{
				Name: "azure-servicebus",
				Kind: "azure.servicebus",
				Properties: map[string]string{
					"queue_name":        dat.queueName,
					"end_point":         dat.endPoint,
					"shared_access_key": dat.sharedAccessKey,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing shared_access_key",
			cfg: config.Spec{
				Name: "azure-servicebus",
				Kind: "azure.servicebus",
				Properties: map[string]string{
					"queue_name":             dat.queueName,
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
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
				t.Errorf("Init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	middle := &mockMiddleware{}
	middle.Init()
	tests := []struct {
		name       string
		cfg        config.Spec
		wantErr    bool
		middleware middleware.Middleware
		timeToWait time.Duration
	}{
		{
			name: "valid servicebus receive",
			cfg: config.Spec{
				Name: "azure-servicebus",
				Kind: "azure.servicebus",
				Properties: map[string]string{
					"queue_name":             dat.queueName,
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
				},
			},
			middleware: middle,
			timeToWait: time.Duration(25) * time.Second,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer leaktest.Check(t)()
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeToWait)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			err = c.Start(ctx, tt.middleware)
			defer func() {
				_ = c.Stop()
			}()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			time.Sleep(tt.timeToWait + 5)
		})
	}
}
