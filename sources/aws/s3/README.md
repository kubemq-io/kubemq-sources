# Kubemq AWS S3 Source Connector

Kubemq AWS S3 source connector allows services using kubemq server to sync aws s3 objects to remote location.

## Prerequisites
The following are required to run the aws s3 source connector:

- kubemq cluster
- kubemq-sources deployment

## Configuration

S3 source connector configuration properties:

| Properties Key    | Required | Description                     | Example                             |
|:------------------|:---------|:--------------------------------|:------------------------------------|
| aws_key        | yes      | aws key                                    | aws key supplied by aws         |
| aws_secret_key | yes      | aws secret key                             | aws secret key supplied by aws  |
| region         | yes      | region                                     | aws region                      |
| token          | no       | aws token ("default" empty string          | aws token                       |
| folders           | yes      | set list of folders to watch    | "/"          |
| target_type       | yes      | set remote target sync type     | "aws", "gcp", "minio", "filesystem" |
| bucket_name       | yes      | set source bucket               | "bucket"                            |
| concurrency       | no       | set sending concurrency         | "1"                                 |
| scan_interval     | no       | set bucket scan interval in sec | "5"                                 |


Example:

```yaml
bindings:
- name: s3
  source:
    kind: aws.s3
    properties:
      aws_key: "id"
      aws_secret_key: 'json'
      region:  "region"
      token: ""
      folders: 'folder1,folder2/sub1/sub2'
      target_type: filesystem
      bucket_name: bucket
      concurrency: 1
      scan_interval: 5
  target:
    kind: kubemq.queue
    properties:
      address: localhost:50000
      channel: queue.s3
  properties: {}
```
