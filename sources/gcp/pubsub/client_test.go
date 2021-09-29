package pubsub

import (
	"context"
	"fmt"
	"github.com/fortytw2/leaktest"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/middleware"
	"github.com/kubemq-io/kubemq-sources/types"
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
	m.channelName = "event.gcp.pubsub"
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
	projectID    string
	credentials  string
	subscriberID string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/projectID.txt")
	if err != nil {
		return nil, err
	}
	t.projectID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/google_cred.json")
	if err != nil {
		return nil, err
	}
	t.credentials = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/subscriberID.txt")
	if err != nil {
		return nil, err
	}
	t.subscriberID = fmt.Sprintf("%s", dat)
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
				Name: "gcp-pubsub",
				Kind: "gcp.pubsub",
				Properties: map[string]string{
					"project_id":    dat.projectID,
					"subscriber_id": dat.subscriberID,
					"credentials":   dat.credentials,
				},
			},
			wantErr: false,
		}, {
			name: "invalid init-missing-credentials",
			cfg: config.Spec{
				Name: "gcp-pubsub",
				Kind: "gcp.pubsub",
				Properties: map[string]string{
					"project_id":    dat.projectID,
					"subscriber_id": dat.subscriberID,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init-missing-project-id",
			cfg: config.Spec{
				Name: "gcp-pubsub",
				Kind: "gcp.pubsub",
				Properties: map[string]string{
					"credentials":   dat.credentials,
					"subscriber_id": dat.subscriberID,
				},
			},
			wantErr: true,
		}, {
			name: "invalid init-missing-subscriber_id",
			cfg: config.Spec{
				Name: "gcp-pubsub",
				Kind: "gcp.pubsub",
				Properties: map[string]string{
					"credentials": dat.credentials,
					"project_id":  dat.projectID,
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

			err := c.Init(ctx, tt.cfg, nil)
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
				Name: "gcp-pubsub",
				Kind: "gcp.pubsub",
				Properties: map[string]string{
					"project_id":    dat.projectID,
					"subscriber_id": dat.subscriberID,
					"credentials":   dat.credentials,
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
			err := c.Init(ctx, tt.cfg, nil)
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
