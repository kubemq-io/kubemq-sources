package sqs

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
	r.SetMetadataKeyValue("result", "ok")
	return r, nil
}

type testStructure struct {
	awsKey       string
	awsSecretKey string
	region       string
	token        string

	sqsQueue   string
	deadLetter string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/aws/awsKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/awsSecretKey.txt")
	if err != nil {
		return nil, err
	}
	t.awsSecretKey = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/region.txt")
	if err != nil {
		return nil, err
	}
	t.region = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/aws/sqs/queue.txt")
	if err != nil {
		return nil, err
	}
	t.sqsQueue = string(dat)

	dat, err = ioutil.ReadFile("./../../../credentials/aws/sqs/deadLetter.txt")
	if err != nil {
		return nil, err
	}
	t.deadLetter = string(dat)
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
				Name: "source-aws-sqs",
				Kind: "source.aws.sqs",
				Properties: map[string]string{
					"aws_key":                dat.awsKey,
					"aws_secret_key":         dat.awsSecretKey,
					"token":                  dat.token,
					"region":                 dat.region,
					"max_number_of_messages": "10",
					"pull_delay":             "15",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error no region",
			cfg: config.Spec{
				Name: "source-aws-sqs",
				Kind: "source.aws.sqs",
				Properties: map[string]string{
					"aws_key":                dat.awsKey,
					"aws_secret_key":         dat.awsSecretKey,
					"token":                  dat.token,
					"max_number_of_messages": "10",
					"pull_delay":             "15",
				},
			},
			wantErr: true,
		}, {
			name: "init - error no aws_key",
			cfg: config.Spec{
				Name: "source-aws-sqs",
				Kind: "source.aws.sqs",
				Properties: map[string]string{
					"aws_secret_key":         dat.awsSecretKey,
					"token":                  dat.token,
					"region":                 dat.region,
					"max_number_of_messages": "10",
					"pull_delay":             "15",
				},
			},
			wantErr: true,
		},
		{
			name: "init -error no aws_secret_key",
			cfg: config.Spec{
				Name: "source-aws-sqs",
				Kind: "source.aws.sqs",
				Properties: map[string]string{
					"aws_key":                dat.awsKey,
					"token":                  dat.token,
					"region":                 dat.region,
					"max_number_of_messages": "1",
					"pull_delay":             "15",
					"visibility_timeout":     "10",
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
			require.EqualValues(t, tt.cfg.Name, c.Name())
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
	}{
		{
			name: "valid sqs receive",
			cfg: config.Spec{
				Name: "source-aws-sqs",
				Kind: "source.aws.sqs",
				Properties: map[string]string{
					"aws_key":                dat.awsKey,
					"aws_secret_key":         dat.awsSecretKey,
					"token":                  dat.token,
					"region":                 dat.region,
					"max_number_of_messages": "1",
					"pull_delay":             "2",
					"queue":                  dat.sqsQueue,
					"visibility_timeout":     "10",
				},
			},
			middleware: &mockMiddleware{},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer leaktest.Check(t)()
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
			time.Sleep(time.Duration(15) * time.Second)
			defer cancel()
			require.NoError(t, err)
		})
	}
}
