package filesystem

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
)

var bucketTypeMap = map[string]string{
	"aws":        "aws",
	"gcp":        "gcp",
	"minio":      "minio",
	"hdfs":       "hdfs",
	"azure":      "azure",
	"filesystem": "filesystem",
}

type options struct {
	folders      []string
	concurrency  int
	bucketType   string
	bucketName   string
	backupFolder string
	scanInterval int
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
	o.scanInterval, err = cfg.Properties.ParseIntWithRange("scan_interval", 5, 1, 3600*365)
	if err != nil {
		return options{}, fmt.Errorf("error parsing scan_interval")
	}
	return o, nil
}
