# Kubemq Minio/S3 Source Connector

Kubemq Minio/S3 source connector allows services using kubemq server to sync minio objects to remote location.

## Prerequisites
The following are required to run the minio source connector:

- kubemq cluster
- minio cluster / AWS s3 service
- kubemq-sources deployment

## Configuration

Minio source connector configuration properties:

| Properties Key    | Required | Description                     | Example                             |
|:------------------|:---------|:--------------------------------|:------------------------------------|
| endpoint          | yes      | minio host address              | "localhost:9000"                    |
| use_ssl           | no       | set connection ssl              | "true"                              |
| access_key_id     | yes      | set access key id               | "minio"                             |
| secret_access_key | yes      | set secret access key           | "minio123"                          |
| folders           | yes      | set list of folders to watch    | "/"          |
| target_type       | yes      | set remote target sync type     | "aws", "gcp", "minio", "filesystem","hdfs","azure","pass-through" |
| bucket_name       | yes      | set source bucket               | "bucket"                            |
| concurrency       | no       | set sending concurrency         | "1"                                 |
| scan_interval     | no       | set bucket scan interval in sec | "5"                                 |


Example:

```yaml
bindings:
- name: minio
  source:
    kind: storage.minio
    properties:
      endpoint: "localhost:9000"
      use_ssl: "false"
      access_key_id: "minio"
      secret_access_key: "minio123"
      folders: 'folder1,folder2/sub1/sub2'
      target_type: filesystem
      bucket_name: bucket
      concurrency: 1
      scan_interval: 5
  target:
    kind: kubemq.queue
    properties:
      address: localhost:50000
      channel: queue.minio
  properties: {}
```
