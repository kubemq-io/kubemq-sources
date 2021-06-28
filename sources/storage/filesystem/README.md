# Kubemq Filesystem Source

Kubemq Filesystem source sends .

## Prerequisites
The following are required to run events source connector:

- kubemq cluster
- kubemq-sources deployment


## Configuration

Filesystem source connector configuration properties:

| Properties Key   | Required | Description                         | Example                                    |
|:-----------------|:---------|:------------------------------------|:-------------------------------------------|
| folders              | yes      | set list of folders to watch   | "c:/folder1,c:/folder2" |
| bucket_type            | yes      | set remote target bucket type              | "aws", "gcp", "minio", "filesystem","hdfs","azure","pass-through" |
| bucket_name          | yes      | set remote target bucket/dir name    | "dir1"          |
| concurrency         | no      | set sending concurrency       | "1"                                 |
| backup_folder         | no      | set backup folder for files after send     | "1"                                 |
| scan_interval     | no       | set filesystem scan interval in sec | "5"                                 |

Example:

```yaml
bindings:
- name: fs
  source:
    kind: storage.filesystem
    properties:
      folders: 'd:\test\source,d:\test\folder2'
      bucket_type: aws
      bucket_name: aws_bucket_name
      concurrency: 1
      backup_folder: 'd:\backup'
      scan_interval: 5
  target:
    kind: kubemq.queue
    properties:
      address: localhost:50000
      channel: queue.fs
  properties: {}

```
