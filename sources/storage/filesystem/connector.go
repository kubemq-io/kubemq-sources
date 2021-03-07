package filesystem

import "github.com/kubemq-hub/builder/connector/common"

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("storage.filesystem").
		SetDescription("Local filesystem properties").
		SetName("File System").
		SetProvider("").
		SetCategory("Storage").
		SetTags("s3","minio").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("folders").
				SetDescription("Set local folders directory to scan").
				SetMust(true),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("bucket_type").
				SetOptions([]string{"aws", "gcp", "minio", "filesystem"}).
				SetDescription("Set remote target bucket type").
				SetMust(true).
				SetDefault("aws"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("bucket_name").
				SetDescription("Set remote target bucket/dir name").
				SetMust(true).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("concurrency").
				SetDescription("Set execution concurrency").
				SetMust(false).
				SetDefault("1"),
		)
}
