package filesystem

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

var bucketTypeMap = map[string]string{
	"aws":        "aws",
	"gcp":        "gcp",
	"minio":      "minio",
	"filesystem": "filesystem",
}

type options struct {
	folders      []string
	concurrency  int
	bucketType   string
	bucketName   string
	backupFolder string
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.folders, err = cfg.Properties.MustParseStringList("folders")
	if err != nil {
		return options{}, fmt.Errorf("error parsing folders, %w", err)
	}
	o.bucketType, err = cfg.Properties.ParseStringMap("bucket_type", bucketTypeMap)
	if err != nil {
		return options{}, fmt.Errorf("error parsing bucket_type, %w", err)
	}
	o.bucketName, err = cfg.Properties.MustParseString("bucket_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing bucket_name, %w", err)
	}
	o.backupFolder = cfg.Properties.ParseString("backup_folder", "")
	o.concurrency, err = cfg.Properties.ParseIntWithRange("concurrency", 1, 1, 1024)
	if err != nil {
		return options{}, fmt.Errorf("error parsing concurrency")
	}
	return o, nil
}
