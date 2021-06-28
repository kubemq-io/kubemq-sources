package s3

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("aws.s3").
		SetDescription("AWS S3 Source").
		SetName("S3").
		SetProvider("AWS").
		SetCategory("Storage").
		SetTags("s3").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_key").
				SetDescription("Set S3 aws key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("aws_secret_key").
				SetDescription("Set S3 aws secret key").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("region").
				SetDescription("Set S3 aws region").
				SetMust(true).
				SetDefault(""),
		).AddProperty(
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
				SetKind("string").
				SetName("token").
				SetDescription("Set S3 token").
				SetMust(false).
				SetDefault(""),
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
