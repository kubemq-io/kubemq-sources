package filesystem

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubemq-io/kubemq-sources/types"
)

type SourceFile struct {
	Info     os.FileInfo
	Path     string
	Root     string
	MovePath string
}

func NewSourceFile(info os.FileInfo, path string, root string, movePath string) *SourceFile {
	return &SourceFile{
		Info:     info,
		Path:     path,
		Root:     root,
		MovePath: movePath,
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

func (s *SourceFile) Metadata() string {
	return fmt.Sprintf("file: %s, size: %d bytes", s.FullPath(), s.Info.Size())
}

func (s *SourceFile) Load() ([]byte, error) {
	return ioutil.ReadFile(s.FullPath())
}

func (s *SourceFile) Hash() string {
	return fmt.Sprintf("%d", s.Info.ModTime().UnixNano())
}

func (s *SourceFile) Do() error {
	if s.MovePath != "" {
		newFileName := strings.Replace(s.FullPath(), s.Root, s.MovePath, 1)
		if err := movefile(s.FullPath(), newFileName); err != nil {
			return err
		}
		return nil
	} else {
		return os.Remove(s.FullPath())
	}
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

func movefile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat error: %s", err)
	}

	err = os.MkdirAll(filepath.Dir(dst), si.Mode())
	if err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		_ = in.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	_ = in.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("sync error: %s", err)
	}

	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("chmod error: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}
