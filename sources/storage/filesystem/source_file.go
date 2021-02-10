package filesystem

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SourceFile struct {
	Info os.FileInfo
	Path string
	Root string
}

func NewSourceFile(info os.FileInfo, path string, root string) *SourceFile {
	return &SourceFile{
		Info: info,
		Path: path,
		Root: root,
	}
}
func (s *SourceFile) FullPath() string {
	p, _ := filepath.Abs(s.Path)
	return filepath.Clean(p)
}
func (s *SourceFile) FileDir() string {
	dir, _ := filepath.Split(s.Path)
	fileDir := strings.Replace(filepath.Clean(dir), filepath.Clean(s.Root), "", -1)
	return fileDir
}
func (s *SourceFile) FileName() string {
	return s.Info.Name()
}
func (s *SourceFile) Load() ([]byte, error) {
	return ioutil.ReadFile(s.FullPath())
}

func (s *SourceFile) Delete() error {
	return os.Remove(s.FullPath())
}
func (s *SourceFile) Request(bucketType string, bucketName string) (*types.Request, error) {
	data, err := s.Load()
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
	case "aws":
		unixFileName := strings.Replace(filepath.Join(s.FileDir(), s.FileName()), `\`, "/", -1)
		unixFileName = strings.TrimPrefix(unixFileName, "/")
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "upload_item").
			SetMetadataKeyValue("bucket_name", bucketName).
			SetMetadataKeyValue("item_name", filepath.Join(s.FileDir(), s.FileName())).
			SetData(data)
	case "minio":
		unixFileName := strings.Replace(filepath.Join(s.FileDir(), s.FileName()), `\`, "/", -1)
		unixFileName = strings.TrimPrefix(unixFileName, "/")
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "put").
			SetMetadataKeyValue("param1", bucketName).
			SetMetadataKeyValue("param2", unixFileName).
			SetData(data)
	case "filesystem":
		targetRequest = NewTargetsRequest().
			SetMetadataKeyValue("method", "save").
			SetMetadataKeyValue("path", filepath.Join(bucketName, s.FileDir())).
			SetMetadataKeyValue("filename", s.FileName()).
			SetData(data)
	default:
		return nil, fmt.Errorf("invalid bucket type")
	}

	return types.NewRequest().SetData(targetRequest.MarshalBinary()), nil
}
