package s3

import (
	"fmt"
	"strings"

	"github.com/kubemq-io/kubemq-sources/config"
)

var bucketTypeMap = map[string]string{
	"aws":          "aws",
	"gcp":          "gcp",
	"minio":        "minio",
	"filesystem":   "filesystem",
	"hdfs":         "hdfs",
	"azure":        "azure",
	"pass-through": "pass-through",
}

type options struct {
	awsKey       string
	awsSecretKey string
	region       string
	token        string
	folders      []string
	concurrency  int
	targetType   string
	bucketName   string
	scanInterval int
}

func parseOptions(cfg config.Spec) (options, error) {
	o := options{}
	var err error
	o.awsKey, err = cfg.Properties.MustParseString("aws_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing aws_key , %w", err)
	}

	o.awsSecretKey, err = cfg.Properties.MustParseString("aws_secret_key")
	if err != nil {
		return options{}, fmt.Errorf("error parsing aws_secret_key , %w", err)
	}

	o.region, err = cfg.Properties.MustParseString("region")
	if err != nil {
		return options{}, fmt.Errorf("error parsing region , %w", err)
	}

	o.token = cfg.Properties.ParseString("token", "")
	o.folders, err = cfg.Properties.MustParseStringList("folders")
	if err != nil {
		return options{}, fmt.Errorf("error parsing folders, %w", err)
	}
	o.targetType, err = cfg.Properties.ParseStringMap("target_type", bucketTypeMap)
	if err != nil {
		return options{}, fmt.Errorf("error parsing target_type, %w", err)
	}
	o.bucketName, err = cfg.Properties.MustParseString("bucket_name")
	if err != nil {
		return options{}, fmt.Errorf("error parsing bucket_name, %w", err)
	}
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

func unixNormalize(in string) string {
	return strings.Replace(in, `\`, "/", -1)
}
