package minio

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-sources/types"
	"github.com/minio/minio-go/v7"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type SourceFile struct {
	Object minio.ObjectInfo
	Bucket string
	client *minio.Client
}

func NewSourceFile(c *minio.Client, bucket string, obj minio.ObjectInfo) *SourceFile {
	return &SourceFile{
		Object: obj,
		Bucket: bucket,
		client: c,
	}
}

func (s *SourceFile) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Bucket, s.Object.Key)
}
func (s *SourceFile) FileDir() string {
	parts := strings.Split(s.Object.Key, "/")
	if len(parts) < 2 {
		return ""
	}
	if len(parts) == 2 {
		return parts[0]
	}
	return strings.Join(parts[:len(parts)-1], "/")
}
func (s *SourceFile) RootDir() string {
	parts := strings.Split(s.Object.Key, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}
func (s *SourceFile) FileName() string {
	parts := strings.Split(s.Object.Key, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
func (s *SourceFile) Metadata() string {
	return fmt.Sprintf("file: %s, size: %d bytes", s.FullPath(), s.Object.Size)
}
func (s *SourceFile) Hash() string {
	return s.Object.ETag

}
func (s *SourceFile) Load(ctx context.Context) ([]byte, error) {
	object, err := s.client.GetObject(ctx, s.Bucket, s.Object.Key, minio.GetObjectOptions{})
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

func (s *SourceFile) Do(ctx context.Context) error {
	return s.Delete(ctx)
}
func (s *SourceFile) Delete(ctx context.Context) error {
	err := s.client.RemoveObject(ctx, s.Bucket, s.Object.Key, minio.RemoveObjectOptions{
		GovernanceBypass: false,
		VersionID:        "",
		Internal:         minio.AdvancedRemoveOptions{},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *SourceFile) Request(ctx context.Context, bucketType string, bucketName string) (*types.Request, error) {
	data, err := s.Load(ctx)
	if err != nil {
		return nil, err
	}
	var targetRequest *TargetsRequest
	switch bucketType {
	case "gcp":
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "upload").
			SetMetadataKeyValue("bucket", bucketName).
			SetMetadataKeyValue("object", s.FileName()).
			SetMetadataKeyValue("path", strings.Replace(s.FileDir(), `\`, "/", -1)).
			SetData(data)
	case "pass-through":
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("path", s.FileDir()).
			SetMetadataKeyValue("filename", s.FileName()).
			SetData(data)
	case "aws":
		unixFileName := strings.Replace(filepath.Join(s.FileDir(), s.FileName()), `\`, "/", -1)
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "upload_item").
			SetMetadataKeyValue("bucket_name", bucketName).
			SetMetadataKeyValue("item_name", strings.TrimPrefix(unixFileName, "/")).
			SetData(data)
	case "minio":
		unixFileName := strings.Replace(filepath.Join(s.FileDir(), s.FileName()), `\`, "/", -1)
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "put").
			SetMetadataKeyValue("param1", bucketName).
			SetMetadataKeyValue("param2", strings.TrimPrefix(unixFileName, "/")).
			SetData(data)
	case "filesystem":
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "save").
			SetMetadataKeyValue("path", s.FileDir()).
			SetMetadataKeyValue("filename", s.FileName()).
			SetData(data)
	case "hdfs":
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "file_path").
			SetMetadataKeyValue("file_path", strings.Replace(s.FullPath(), `\`, "/", -1)).
			SetData(data)
	case "azure":
		unixFileName := strings.Replace(filepath.Join(s.FileDir(), s.FileName()), `\`, "/", -1)
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "upload").
			SetMetadataKeyValue("service_url", strings.TrimPrefix(unixFileName, "/")).
			SetData(data)
	default:
		return nil, fmt.Errorf("invalid target type")
	}

	return types.NewRequest().SetData(targetRequest.MarshalBinary()), nil
}
