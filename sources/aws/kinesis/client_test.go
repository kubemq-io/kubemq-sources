package kinesis

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

	ShardIteratorType string
	consumerARN       string
	shardID           string
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

	t.ShardIteratorType = "LATEST"

	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/consumerARN.txt")
	if err != nil {
		return nil, err
	}
	t.consumerARN = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/aws/kinesis/shardID.txt")
	if err != nil {
		return nil, err
	}
	t.shardID = string(dat)
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
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_key":             dat.awsKey,
					"aws_secret_key":      dat.awsSecretKey,
					"token":               dat.token,
					"region":              dat.region,
					"shard_id":            dat.shardID,
					"consumer_arn":        dat.consumerARN,
					"shard_iterator_type": dat.ShardIteratorType,
				},
			},
			wantErr: false,
		},
		{
			name: "init - error no region",
			cfg: config.Spec{
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_key":             dat.awsKey,
					"aws_secret_key":      dat.awsSecretKey,
					"token":               dat.token,
					"shard_id":            dat.shardID,
					"consumer_arn":        dat.consumerARN,
					"shard_iterator_type": dat.ShardIteratorType,
				},
			},
			wantErr: true,
		}, {
			name: "init - error no aws_key",
			cfg: config.Spec{
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_secret_key":      dat.awsSecretKey,
					"token":               dat.token,
					"region":              dat.region,
					"shard_id":            dat.shardID,
					"consumer_arn":        dat.consumerARN,
					"shard_iterator_type": dat.ShardIteratorType,
				},
			},
			wantErr: true,
		},
		{
			name: "init -error no aws_secret_key",
			cfg: config.Spec{
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_key":             dat.awsKey,
					"token":               dat.token,
					"region":              dat.region,
					"shard_id":            dat.shardID,
					"consumer_arn":        dat.consumerARN,
					"shard_iterator_type": dat.ShardIteratorType,
				},
			},
			wantErr: true,
		}, {
			name: "init -error no shard_iterator_type",
			cfg: config.Spec{
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_key":        dat.awsKey,
					"aws_secret_key": dat.awsSecretKey,
					"token":          dat.token,
					"region":         dat.region,
					"shard_id":       dat.shardID,
					"consumer_arn":   dat.consumerARN,
				},
			},
			wantErr: true,
		}, {
			name: "init - error no consumer_arn",
			cfg: config.Spec{
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_key":             dat.awsKey,
					"aws_secret_key":      dat.awsSecretKey,
					"token":               dat.token,
					"region":              dat.region,
					"shard_id":            dat.shardID,
					"shard_iterator_type": dat.ShardIteratorType,
				},
			},
			wantErr: true,
		},
		{
			name: "init -error no shard_id",
			cfg: config.Spec{
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_key":             dat.awsKey,
					"aws_secret_key":      dat.awsSecretKey,
					"token":               dat.token,
					"region":              dat.region,
					"shard_iterator_type": dat.ShardIteratorType,
					"consumer_arn":        dat.consumerARN,
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
			name: "valid kinesis receive",
			cfg: config.Spec{
				Name: "source-aws-kinesis",
				Kind: "source.aws.kinesis",
				Properties: map[string]string{
					"aws_key":             dat.awsKey,
					"aws_secret_key":      dat.awsSecretKey,
					"token":               dat.token,
					"region":              dat.region,
					"shard_id":            dat.shardID,
					"consumer_arn":        dat.consumerARN,
					"shard_iterator_type": dat.ShardIteratorType,
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
