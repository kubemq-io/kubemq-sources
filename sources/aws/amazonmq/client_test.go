package amazonmq

import (
	"context"
	"fmt"
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

type testStructure struct {
	host        string
	username    string
	password    string
	destination string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/aws/amazonmq/host.txt")
	if err != nil {
		return nil, err
	}
	t.host = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/amazonmq/username.txt")
	if err != nil {
		return nil, err
	}
	t.username = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/amazonmq/password.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/amazonmq/destination.txt")
	if err != nil {
		return nil, err
	}
	t.destination = string(dat)
	return t, nil
}

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
	m.channelName = "event.aws.amazonmq"
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
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":            dat.host,
					"username":        dat.username,
					"password":        dat.password,
					"destination":     dat.destination,
					"dynamic_mapping": "false",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - bad url",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":            "localhost:41231",
					"username":        dat.username,
					"password":        dat.password,
					"destination":     dat.destination,
					"dynamic_mapping": "false",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - no  url",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":            "fake",
					"username":        dat.username,
					"password":        dat.password,
					"destination":     dat.destination,
					"dynamic_mapping": "false",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - no destination",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":            dat.host,
					"username":        dat.username,
					"password":        dat.password,
					"dynamic_mapping": "false",
				},
			},
			wantErr: true,
		}, {
			name: "invalid init - missing dynamic_mapping",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":        dat.host,
					"username":    dat.username,
					"password":    dat.password,
					"destination": dat.destination,
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
	dat, err := getTestStructure()
	require.NoError(t, err)
	middle := &mockMiddleware{}
	middle.Init()
	tests := []struct {
		name       string
		cfg        config.Spec
		wantErr    bool
		middleware middleware.Middleware
	}{
		{
			name: "valid amazonmq receive",
			cfg: config.Spec{
				Name: "aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":            dat.host,
					"username":        dat.username,
					"password":        dat.password,
					"destination":     dat.destination,
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
			time.Sleep(time.Duration(5) * time.Second)
			require.NoError(t, err)
		})
	}
}
