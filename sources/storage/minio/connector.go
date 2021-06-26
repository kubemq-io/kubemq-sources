package minio

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("storage.minio").
		SetDescription("Minio Storage Source").
		SetName("Minio").
		SetProvider("").
		SetCategory("Storage").
		SetTags("s3").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("endpoint").
				SetTitle("Endpoint").
				SetDescription("Set Minio endpoint address").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("use_ssl").
				SetTitle("USE SSL").
				SetDescription("Set Minio SSL connection").
				SetMust(false).
				SetDefault("true"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("access_key_id").
				SetTitle("Access Key ID").
				SetDescription("Set Minio access key id").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("multilines").
				SetName("secret_access_key").
				SetTitle("Access Key Secret").
				SetDescription("Set Minio secret access key").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("bucket_name").
				SetTitle("Bucket Source").
				SetDescription("Set remote target bucket/dir name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("folders").
				SetTitle("Sync Folders Names").
				SetDescription("Set bucket folders directory to scan").
				SetMust(true).
				SetDefault("/"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("target_type").
				SetTitle("Sync Target Type").
				SetOptions([]string{"aws", "gcp", "minio", "filesystem", "hdfs", "azure", "pass-through"}).
				SetDescription("Set remote target bucket type").
				SetMust(true).
				SetDefault("filesystem"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("concurrency").
				SetDescription("Set execution concurrency").
				SetMust(false).
				SetDefault("1"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("scan_interval").
				SetDescription("Set scan interval in seconds").
				SetMust(false).
				SetDefault("5"),
		)

}
