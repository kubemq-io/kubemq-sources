package pubsub

import (
	"context"
	"fmt"
	"github.com/fortytw2/leaktest"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type mockMiddleware struct {
}

func (m *mockMiddleware) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	fmt.Println(request)
	r := types.NewResponse()
	r.SetData([]byte("ok"))
	r.SetMetadata(`"result":"ok"`)
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
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"project_id":    dat.projectID,
					"subscriber_id": dat.subscriberID,
					"credentials":   dat.credentials,
				},
			},
			wantErr: false,
		}, {
			name: "init-missing-credentials",
			cfg: config.Spec{
				Name: "target-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"project_id":    dat.projectID,
					"subscriber_id": dat.subscriberID,
				},
			},
			wantErr: true,
		},
		{
			name: "init-missing-project-id",
			cfg: config.Spec{
				Name: "source-gcp-pubsub",
				Kind: "target.gcp.pubsub",
				Properties: map[string]string{
					"credentials":   dat.credentials,
					"subscriber_id": dat.subscriberID,
				},
			},
			wantErr: true,
		}, {
			name: "init-missing-subscriber_id",
			cfg: config.Spec{
				Name: "source-gcp-pubsub",
				Kind: "target.gcp.pubsub",
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
				Name: "source-gcp-pubsub",
				Kind: "source.gcp.pubsub",
				Properties: map[string]string{
					"project_id":    dat.projectID,
					"subscriber_id": dat.subscriberID,
					"credentials":   dat.credentials,
				},
			},
			middleware: &mockMiddleware{},
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
			time.Sleep(tt.timeToWait + 5)
			err = c.Stop()
			require.NoError(t, err)
			time.Sleep(tt.timeToWait + 2)
			defer cancel()
			require.NoError(t, err)
		})
	}
}
