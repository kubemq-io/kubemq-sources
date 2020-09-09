# Kubemq kinesis source Connector

Kubemq aws-kinesis source connector allows services using kubemq server to access aws kinesis service.

## Prerequisites
The following required to run the aws-kinesis source connector:

- kubemq cluster
- aws account with kinesis active service
- kubemq-source deployment

## Configuration

kinesis source connector configuration properties:

| Properties Key                 | Required | Description                                                       | Example                     |
|:-------------------------------|:---------|:------------------------------------------------------------------|:----------------------------|
| aws_key                        | yes      | aws key                                                           | aws key supplied by aws         |
| aws_secret_key                 | yes      | aws secret key                                                    | aws secret key supplied by aws  |
| region                         | yes      | region                                                            | aws region                      |
| consumer_arn                   | yes      | consumer_arn generates when creating consumer                     | "arn:myid"                      |
| shard_iterator_type            | yes      | Determines how the shard iterator is used to start                | "AT_SEQUENCE_NUMBER","AFTER_SEQUENCE_NUMBER","TRIM_HORIZON","LATEST"|
| sequence_number                | no       | sequence to start streaming at (default "")                       | "1" (default 1)                 |
| shard_id                       | yes      | The unique identifier of the shard                                | "my_shard"                      |
| pull_delay                     | no       | wait time between calls (milliseconds)                            | "1" (default 5)                 |
 

Example:

```yaml
bindings:
  - name: kubemq-query-aws-kinesis
    source:
      kind: source.query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-aws-kinesis-connector"
        auth_token: ""
        channel: "query.aws.kinesis"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: source.aws.kinesis
      name: source-aws-kinesis
      properties:
        aws_key: "id"
        aws_secret_key: 'json'
        region:  "my region"
        consumer_arn : "arn:my_consumer"
        shard_iterator_type : "LATEST"
        sequence_number : "2341"
        shard_id : "my_shard-123456"
        pull_delay : "5"
```


