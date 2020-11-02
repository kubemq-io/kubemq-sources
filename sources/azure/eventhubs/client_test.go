package eventhubs

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
	partitionID         string
	endPoint            string
	sharedAccessKeyName string
	sharedAccessKey     string
	entityPath          string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/azure/eventhubs/partitionID.txt")
	if err != nil {
		return nil, err
	}
	t.partitionID = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/eventhubs/endPoint.txt")
	if err != nil {
		return nil, err
	}
	t.endPoint = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/eventhubs/sharedAccessKeyName.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKeyName = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/eventhubs/sharedAccessKey.txt")
	if err != nil {
		return nil, err
	}
	t.sharedAccessKey = fmt.Sprintf("%s", dat)
	dat, err = ioutil.ReadFile("./../../../credentials/azure/eventhubs/entityPath.txt")
	if err != nil {
		return nil, err
	}
	t.entityPath = fmt.Sprintf("%s", dat)
	
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
				Name: "aws-eventhubs",
				Kind: "aws.eventhubs",
				Properties: map[string]string{
					"partition_id":           dat.partitionID,
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
					"entity_path":            dat.entityPath,
				},
			},
			wantErr: false,
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
			name: "valid eventhubs receive",
			cfg: config.Spec{
				Name: "aws-eventhubs",
				Kind: "aws.eventhubs",
				Properties: map[string]string{
					"partition_id":           dat.partitionID,
					"end_point":              dat.endPoint,
					"shared_access_key_name": dat.sharedAccessKeyName,
					"shared_access_key":      dat.sharedAccessKey,
					"entity_path":            dat.entityPath,
				},
			},
			middleware: &mockMiddleware{},
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
				c.Stop()
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
