package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kubemq-hub/kubemq-sources/types"

	"github.com/google/uuid"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type mockTarget struct {
	requests []*TargetsRequest
}

func (m *mockTarget) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	tr := &TargetsRequest{}
	err := json.Unmarshal(request.Data, tr)
	if err != nil {
		return nil, err
	}
	m.requests = append(m.requests, tr)
	return types.NewResponse(), nil
}

type minioTestClient struct {
	client *minio.Client
}

func (c *minioTestClient) MakeBucket(ctx context.Context, name string) error {
	bucketOptions := minio.MakeBucketOptions{

		ObjectLocking: false,
	}
	err := c.client.MakeBucket(ctx, name, bucketOptions)
	if err != nil {
		return err
	}
	return nil
}

func (c *minioTestClient) RemoveBucket(ctx context.Context, name string) error {
	err := c.client.RemoveBucket(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

func (c *minioTestClient) ListObjects(ctx context.Context, bucket string) ([]minio.ObjectInfo, error) {
	var objects []minio.ObjectInfo
	for object := range c.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
		objects = append(objects, object)
	}
	return objects, nil

}
func (c *minioTestClient) Get(ctx context.Context, bucket, name string) ([]byte, error) {
	object, err := c.client.GetObject(ctx, bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = object.Close()
	}()
	data, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (c *minioTestClient) Put(ctx context.Context, bucket, name string, value []byte) error {
	r := bytes.NewReader(value)
	_, err := c.client.PutObject(ctx, bucket, name, r, int64(r.Len()), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *minioTestClient) Remove(ctx context.Context, bucket, name string) error {
	err := c.client.RemoveObject(ctx, bucket, name, minio.RemoveObjectOptions{
		GovernanceBypass: false,
		VersionID:        "",
		Internal:         minio.AdvancedRemoveOptions{},
	})
	if err != nil {
		return err
	}
	return nil
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
				Name: "minio-source",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "localhost:9000",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
					"bucket_name":       "bucket2",
					"folders":           "folder1,folder2",
					"target_type":       "filesystem",
				},
			},
			wantErr: false,
		},
		{
			name: "init - no endpoint key",
			cfg: config.Spec{
				Name: "minio-source",
				Kind: "",
				Properties: map[string]string{
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad endpoint key",
			cfg: config.Spec{
				Name: "minio-source",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "badhost",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no access key",
			cfg: config.Spec{
				Name: "minio-source",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "localhost:9001",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no secret key",
			cfg: config.Spec{
				Name: "minio-source",
				Kind: "",
				Properties: map[string]string{
					"endpoint":      "localhost:9001",
					"access_key_id": "minio",
					"use_ssl":       "false",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad bucket name",
			cfg: config.Spec{
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "localhost:9000",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
					"bucket_name":       "bucketsasdasd",
					"folders":           "folder1,folder2",
					"target_type":       "filesystem",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad folders",
			cfg: config.Spec{
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "localhost:9000",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
					"bucket_name":       "bucket",
					"folders":           "",
					"target_type":       "filesystem",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad bucket type",
			cfg: config.Spec{
				Name: "minio-target",
				Kind: "",
				Properties: map[string]string{
					"endpoint":          "localhost:9000",
					"access_key_id":     "minio",
					"secret_access_key": "minio123",
					"use_ssl":           "false",
					"bucket_name":       "bucket",
					"folders":           "folder1,folder2",
					"target_type":       "",
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

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantPutErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
func TestClient_Flow(t *testing.T) {
	bucket := uuid.New().String()
	cfg := config.Spec{
		Name: "minio-source",
		Kind: "",
		Properties: map[string]string{
			"endpoint":          "localhost:9000",
			"access_key_id":     "minio",
			"secret_access_key": "minio123",
			"use_ssl":           "false",
			"bucket_name":       bucket,
			"folders":           "folder1,folder2",
			"target_type":       "filesystem",
			"scan_interval":     "1",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := New()
	err := c.Init(ctx, cfg, nil)
	require.NoError(t, err)
	testClient := &minioTestClient{
		client: c.s3Client,
	}
	err = testClient.MakeBucket(ctx, bucket)
	require.NoError(t, err)

	mock := &mockTarget{
		requests: nil,
	}

	err = c.Start(ctx, mock)
	err = testClient.Put(ctx, bucket, "folder1/file1.txt", []byte("data"))
	require.NoError(t, err)
	err = testClient.Put(ctx, bucket, "folder2/sub1/sub2/file2.txt", []byte("data"))
	require.NoError(t, err)

	time.Sleep(time.Duration(c.opts.scanInterval*8) * time.Second)

	require.EqualValues(t, 2, len(mock.requests))
	for _, request := range mock.requests {
		require.EqualValues(t, request.Data, []byte("data"))
	}

	_, err = testClient.Get(ctx, bucket, "folder1/file1.txt")
	require.Error(t, err)

	_, err = testClient.Get(ctx, bucket, "folder2/sub1/sub2/file2.txt")
	require.Error(t, err)

	err = testClient.RemoveBucket(ctx, bucket)
	require.NoError(t, err)
}
