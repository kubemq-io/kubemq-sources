package amazonmq

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
				Name: "source-aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":        dat.host,
					"username":    dat.username,
					"password":    dat.password,
					"destination": dat.destination,
				},
			},
			wantErr: false,
		}, {
			name: "init - no host",
			cfg: config.Spec{
				Name: "source-aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"username":    dat.username,
					"password":    dat.password,
					"destination": dat.destination,
				},
			},
			wantErr: true,
		}, {
			name: "init - no destination",
			cfg: config.Spec{
				Name: "source-aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"username": dat.username,
					"password": dat.password,
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
	tests := []struct {
		name       string
		cfg        config.Spec
		wantErr    bool
		middleware middleware.Middleware
		timeToWait time.Duration
	}{
		{
			name: "valid amazonmq receive",
			cfg: config.Spec{
				Name: "source-aws-amazonmq",
				Kind: "aws.amazonmq",
				Properties: map[string]string{
					"host":        dat.host,
					"username":    dat.username,
					"password":    dat.password,
					"destination": dat.destination,
				},
			},
			middleware: &mockMiddleware{},
			timeToWait: time.Duration(5) * time.Second,
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
