// +build container

package ibmmq

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
	m.channelName = "event.messaging.ibmmq"
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
	applicationChannelName string
	hostname               string
	listenerPort           string
	queueManagerName       string
	apiKey                 string
	mqUsername             string
	password               string
	QueueName              string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/applicationChannelName.txt")
	if err != nil {
		return nil, err
	}
	t.applicationChannelName = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/hostname.txt")
	if err != nil {
		return nil, err
	}
	t.hostname = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/listenerPort.txt")
	if err != nil {
		return nil, err
	}
	t.listenerPort = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/connectionInfo/queueManagerName.txt")
	if err != nil {
		return nil, err
	}
	t.queueManagerName = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/applicationApiKey/apiKey.txt")
	if err != nil {
		return nil, err
	}
	t.apiKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/applicationApiKey/mqUsername.txt")
	if err != nil {
		return nil, err
	}
	t.mqUsername = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/ibm/mq/applicationApiKey/mqPassword.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)
	t.QueueName = "DEV.QUEUE.1"
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
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init - missing host_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing queue_manager_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"host_name":      dat.hostname,
					"port_number":    dat.listenerPort,
					"channel_name":   dat.applicationChannelName,
					"username":       dat.mqUsername,
					"key_repository": dat.apiKey,
					"password":       dat.password,
					"queue_name":     dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing channel_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing username",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing queue_name",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
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

			err := c.Init(ctx, tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				t.Logf("init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)

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
			name: "valid pubsub receive",
			cfg: config.Spec{
				Name: "messaging-ibmmq",
				Kind: "messaging.ibmmq",
				Properties: map[string]string{
					"queue_manager_name": dat.queueManagerName,
					"host_name":          dat.hostname,
					"port_number":        dat.listenerPort,
					"channel_name":       dat.applicationChannelName,
					"username":           dat.mqUsername,
					"key_repository":     dat.apiKey,
					"password":           dat.password,
					"queue_name":         dat.QueueName,
				},
			},
			middleware: middle,
			timeToWait: time.Duration(60) * time.Second,
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
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			time.Sleep(tt.timeToWait + 5)
			err = c.Stop()
			require.NoError(t, err)
			time.Sleep(tt.timeToWait + 2)
			require.NoError(t, err)
		})
	}
}
