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
| root              | yes      | set root path for file watching   | "c:/file" |
| bucket_type            | yes      | set remote target bucket type              | "aws", "gcp", "minio", "filesystem" |
| bucket_name          | yes      | set remote target bucket/dir name    | "dir1"          |
| concurrency         | no      | set sending concurrency       | "1"                                 |

Example:

```yaml
bindings:
- name: fs
  source:
    kind: storage.filesystem
    properties:
      root: 'd:\test\source'
      bucket_type: filesystem
      bucket_name: bucket
      concurrency: 5
  target:
    kind: kubemq.queue
    properties:
      address: localhost:50000
      channel: queue.fs
  properties: {}

```
