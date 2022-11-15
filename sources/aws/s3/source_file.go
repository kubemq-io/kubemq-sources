package s3

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kubemq-io/kubemq-sources/types"
)

type SourceFile struct {
	Object     *s3.Object
	Bucket     string
	downloader *s3manager.Downloader
	client     *s3.S3
}

func NewSourceFile(c *s3.S3, downloader *s3manager.Downloader, bucket string, obj *s3.Object) *SourceFile {
	return &SourceFile{
		Object:     obj,
		Bucket:     bucket,
		downloader: downloader,
		client:     c,
	}
}

func (s *SourceFile) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Bucket, *s.Object.Key)
}

func (s *SourceFile) FileDir() string {
	parts := strings.Split(*s.Object.Key, "/")
	if len(parts) < 2 {
		return ""
	}
	if len(parts) == 2 {
		return parts[0]
	}
	return strings.Join(parts[:len(parts)-1], "/")
}

func (s *SourceFile) RootDir() string {
	parts := strings.Split(*s.Object.Key, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

func (s *SourceFile) FileName() string {
	parts := strings.Split(*s.Object.Key, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func (s *SourceFile) Metadata() string {
	return fmt.Sprintf("file: %s, size: %d bytes", s.FullPath(), *s.Object.Size)
}

func (s *SourceFile) Hash() string {
	return *s.Object.ETag
}

func (s *SourceFile) Load(ctx context.Context) ([]byte, error) {
	requestInput := s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    s.Object.Key,
	}
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := s.downloader.DownloadWithContext(ctx, buf, &requestInput)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *SourceFile) Do(ctx context.Context) error {
	return s.Delete(ctx)
}

func (s *SourceFile) Delete(ctx context.Context) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    s.Object.Key,
	})
	if err != nil {
		return err
	}
	err = s.client.WaitUntilObjectNotExistsWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    s.Object.Key,
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
