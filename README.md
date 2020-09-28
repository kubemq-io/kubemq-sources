# KubeMQ Sources

KubeMQ Sources connects external systems and cloud services with KubeMQ message queue broker.

KubeMQ Sources allows us to build a message-based microservices architecture on Kubernetes with minimal efforts and without developing connectivity interfaces between external system such as messaging components (RabbitMQ, Kafka, MQTT) ,REST APIs and KubeMQ message queue broker.
In addition, KubeMQ Sources allows to migrate legacy systems (together with [KubeMQ Targets](https://github.com/kubemq-hub/kubemq-targets)) to KubeMQ based messaging architecture.


**Key Features**:

- **Runs anywhere**  - Kubernetes, Cloud, on-prem, anywhere
- **Stand-alone** - small docker container / binary
- **Single Interface** - One interface all the services
- **API Gateway** - Act as an REST Api gateway
- **Plug-in Architecture** Easy to extend, easy to connect
- **Middleware Supports** - Logs, Metrics, Retries, and Rate Limiters
- **Easy Configuration** - simple yaml file builds your topology

## Concepts

KubeMQ Targets building blocks are:
 - Binding
 - Source
 - Target

### Binding

Binding is a 1:1 connection between Source and Target. Every Binding runs independently.

![binding](.github/assets/binding.jpg)

### Source

Target is an external service that exposes an API allowing to interact and serve his functionalists with other services.

Targets can be Cache systems such as Redis and Memcached, SQL Databases such as Postgres and MySql, and event an HTTP generic Rest interface.

KubeMQ Targets integrate each one of the supported targets and service requests based on the request data.

A list of supported targets is below.

#### Standalone Services

| Category   | Target                                                              | Kind                         | Configuration                                      | Example                                 |
|:-----------|:--------------------------------------------------------------------|:-----------------------------|:---------------------------------------------------|:----------------------------------------|
| Cache      |                                                                     |                              |                                                    |                                         |
|            | [Redis](https://redis.io/)                                          | target.cache.redis           | [Usage](targets/cache/redis)                       | [Example](examples/cache/redis)         |
|            | [Memcached](https://memcached.org/)                                 | target.cache.memcached       | [Usage](targets/cache/memcached)                   | [Example](examples/cache/memcached)    |
| Stores/db  |                                                                     |                              |                                                    |                                         |
|            | [Postgres](https://www.postgresql.org/)                             | target.stores.postgres       | [Usage](targets/stores/postgres)                   | [Example](examples/stores/postgres)     |
|            | [Mysql](https://www.mysql.com/)                                     | target.stores.mysql          | [Usage](targets/stores/mysql)                      | [Example](examples/stores/mysql)        |
|            | [MSSql](https://www.microsoft.com/en-us/sql-server/sql-server-2019) | target.stores.mssql          | [Usage](targets/stores/mssql)                      | [Example](examples/stores/mssql)        |
|            | [MongoDB](https://www.mongodb.com/)                                 | target.stores.mongodb        | [Usage](targets/stores/mongodb)                    | [Example](examples/stores/mongodb)      |
|            | [Elastic Search](https://www.elastic.co/)                           | target.stores.elastic-search | [Usage](targets/stores/elastic)                    | [Example](examples/stores/elastic)      |
|            | [Cassandra](https://cassandra.apache.org/)                          | target.stores.cassandra      | [Usage](targets/stores/cassandra)                  | [Example](examples/stores/cassandra)    |
|            | [Couchbase](https://www.couchbase.com/)                             | target.stores.couchbase      | [Usage](targets/stores/couchbase)                  | [Example](examples/stores/couchbase)    |
| Messaging  |                                                                     |                              |                                                    |                                         |
|            | [Kafka](https://kafka.apache.org/)                                  | target.messaging.kafka       | [Usage](targets/messaging/kafka)                   | [Example](examples/messaging/kafka)     |
|            | [RabbitMQ](https://www.rabbitmq.com/)                               | target.messaging.rabbitmq    | [Usage](targets/messaging/rabbitmq)                | [Example](examples/messaging/rabbitmq)  |
|            | [MQTT](http://mqtt.org/)                                            | target.messaging.mqtt        | [Usage](targets/messaging/mqtt)                    | [Example](examples/messaging/mqtt)      |
|            | [ActiveMQ](http://activemq.apache.org/)                             | target.messaging.activemq    | [Usage](targets/messaging/activemq)                | [Example](examples/messaging/activemq)  |
| Storage    |                                                                     |                              |                                                    |                                         |
|            | [Minio/S3](https://min.io/)                                         | target.storage.minio         | [Usage](targets/storage/minio)                     | [Example](examples/storage/minio)       |
| Serverless |                                                                     |                              |                                                    |                                         |
|            | [OpenFaas](https://www.openfaas.com/)                               | target.serverless.openfaas   | [Usage](targets/serverless/openfass)               | [Example](examples/serverless/openfass) |
| Http       |                                                                     |                              |                                                    |                                         |
|            | Http                                                                | target.http                  | [Usage](targets/http)                              | [Example](examples/http)                |




#### Google Cloud Platform (GCP)

| Category   | Target                                                              | Kind                       | Configuration                              | Example                                       |
|:-----------|:--------------------------------------------------------------------|:---------------------------|:-------------------------------------------|:----------------------------------------------|
| Cache      |                                                                     |                            |                                            |                                               |
|            | [Redis](https://cloud.google.com/memorystore)                       | target.gcp.cache.redis     | [Usage](targets/gcp/memorystore/redis)     | [Example](examples/gcp/memorystore/redis)     |
|            | [Memcached](https://cloud.google.com/memorystore)                   | target.gcp.cache.memcached | [Usage](targets/gcp/memorystore/memcached) | [Example](examples/gcp/memorystore/memcached) |
| Stores/db  |                                                                     |                            |                                            |                                               |
|            | [Postgres](https://cloud.google.com/sql)                            | target.gcp.stores.postgres | [Usage](targets/gcp/sql/postgres)          | [Example](examples/gcp/sql/postgres)           |
|            | [Mysql](https://cloud.google.com/sql)                               | target.gcp.stores.mysql    | [Usage](targets/gcp/sql/mysql)             | [Example](examples/gcp/sql/mysql)              |
|            | [BigQuery](https://cloud.google.com/bigquery)                       | target.gcp.bigquery        | [Usage](targets/gcp/bigquery)              | [Example](examples/gcp/bigquery)               |
|            | [BigTable](https://cloud.google.com/bigtable)                       | target.gcp.bigtable        | [Usage](targets/gcp/bigtable)              | [Example](examples/gcp/bigtable)               |
|            | [Firestore](https://cloud.google.com/firestore)                     | target.gcp.firestore       | [Usage](targets/gcp/firestore)             | [Example](examples/gcp/firestore)              |
|            | [Spanner](https://cloud.google.com/spanner)                         | target.gcp.spanner         | [Usage](targets/gcp/spanner)               | [Example](examples/gcp/spanner)                |
|            | [Firebase](https://firebase.google.com/products/realtime-database/) | target.gcp.firebase        | [Usage](targets/gcp/firebase)              | [Example](examples/gcp/firebase)               |
| Messaging  |                                                                     |                            |                                            |                                               |
|            | [Pub/Sub](https://cloud.google.com/pubsub)                          | target.gcp.pubsub          | [Usage](targets/gcp/pubsub)                | [Example](examples/gcp/pubsub)                 |
| Storage    |                                                                     |                            |                                            |                                               |
|            | [Storage](https://cloud.google.com/storage)                         | target.gcp.storage         | [Usage](targets/gcp/storage)               | [Example](examples/gcp/storage)                |
| Serverless |                                                                     |                            |                                            |                                               |
|            | [Functions](https://cloud.google.com/functions)                     | target.gcp.cloudfunctions  | [Usage](targets/gcp/cloudfunctions)        | [Example](examples/gcp/cloudfunctions)         |
|            |                                                                     |                            |                                            |                                               |



#### Amazon Web Service (AWS)


| Category   | Target                                                        | Kind                               | Configuration                                               | Example                                      |
|:-----------|:--------------------------------------------------------------|:-----------------------------------|:------------------------------------------------------------|:---------------------------------------------|
| Stores/db  |                                                               |                                    |                                                             |                                              |
|            | [Athena](https://docs.aws.amazon.com/athena)                  | target.aws.athena                  | [Usage](targets/aws/athena)                                 | [Example](examples/aws/athena)               |
|            | [DynamoDB](https://aws.amazon.com/dynamodb/)                  | target.aws.dynamodb                | [Usage](targets/aws/dynamodb)                               | [Example](examples/aws/dynamodb)             |
|            | [Elasticsearch](https://aws.amazon.com/elasticsearch-service/)| target.aws.elasticsearch           | [Usage](targets/aws/elasticsearch)                          | [Example](examples/aws/elasticsearch)        |
|            | [KeySpaces](https://docs.aws.amazon.com/keyspaces)            | target.aws.keyspaces               | [Usage](targets/aws/keyspaces)                              | [Example](examples/aws/keyspaces)            |
|            | [MariaDB](https://aws.amazon.com/rds/mariadb/)                | target.aws.rds.mariadb             | [Usage](targets/aws/rds/mariadb)                            | [Example](examples/aws/rds/mariadb)          |
|            | [MSSql](https://aws.amazon.com/rds/sqlserver/)                | target.aws.rds.mssql               | [Usage](targets/aws/rds/mssql)                              | [Example](examples/aws/rds/mssql)            |
|            | [MySQL](https://aws.amazon.com/rds/mysql/)                    | target.aws.rds.mysql               | [Usage](targets/aws/rds/mysql)                              | [Example](examples/aws/rds/mysql)            |       
|            | [Postgres](https://aws.amazon.com/rds/postgresql/)            | target.aws.rds.postgres            | [Usage](targets/aws/rds/postgres)                           | [Example](examples/aws/rds/postgres)         |     
|            | [RedShift](https://aws.amazon.com/redshift/)                  | target.aws.rds.redshift            | [Usage](targets/aws/rds/redshift)                           | [Example](examples/aws/rds/redshift)         |
|            | [RedShiftSVC](https://aws.amazon.com/redshift/)               | target.aws.rds.redshift.service    | [Usage](targets/aws/redshift)                               | [Example](examples/aws/redshift)             |
| Messaging  |                                                               |                                    |                                                             |                                              |
|            | [AmazonMQ](https://aws.amazon.com/amazon-mq/)                 | target.aws.amazonmq                | [Usage](targets/aws/amazonmq)                               | [Example](examples/aws/amazonmq)             |
|            | [msk](https://aws.amazon.com/msk/)                            | target.aws.msk                     | [Usage](targets/aws/msk)                                    | [Example](examples/aws/msk)                  |       
|            | [Kinesis](https://aws.amazon.com/kinesis/)                    | target.aws.kinesis                 | [Usage](targets/aws/kinesis)                                | [Example](examples/aws/kinesis)              |  
|            | [SQS](https://aws.amazon.com/sqs/)                            | target.aws.sqs                     | [Usage](targets/aws/sqs)                                    | [Example](examples/aws/sqs)                  |         
|            | [SNS](https://aws.amazon.com/sns/)                            | target.aws.sns                     | [Usage](targets/aws/sns)                                    | [Example](examples/aws/sns)                  |       
| Storage    |                                                               |                                    |                                                             |                                              |
|            | [s3](https://aws.amazon.com/s3/)                              | target.aws.s3                      | [Usage](targets/aws/s3)                                     | [Example](examples/aws/s3)                   |
| Serverless |                                                               |                                    |                                                             |                                              |
|            | [lambda](https://aws.amazon.com/lambda/)                      | target.aws.lambda                  | [Usage](targets/aws/lambda)                                 | [Example](examples/aws/lambda)               | 
| Other      |                                                               |                                    |                                                             |                                              |
|            | [Cloud Watch](https://aws.amazon.com/cloudwatch/)             | target.aws.cloudwatch.logs         | [Usage](targets/aws/cloudwatch/logs)                        | [Example](examples/aws/cloudwatch/logs)      |
|            | [Cloud Watch](https://aws.amazon.com/cloudwatch/)             | target.aws.cloudwatch.events       | [Usage](targets/aws/cloudwatch/events)                      | [Example](examples/aws/cloudwatch/events)    |
|            | [Cloud Watch](https://aws.amazon.com/cloudwatch/)             | target.aws.cloudwatch.metrics      | [Usage](targets/aws/cloudwatch/metrics)                     | [Example](examples/aws/cloudwatch/metrics)   |
|            |                                                               |                                    |                                                             |                                              |


#### Microsoft Azure

(Work in Progress)

### Target

The target is a KubeMQ connection which send the data from the sources and route them to the appropriate KubeMQ channel for action, and return back a response if needed.

KubeMQ Sources supports all of KubeMQ's messaging patterns: Queue, Events, Events-Store, Command, and Query.


| Type                                                                              | Kind                | Configuration                           |
|:----------------------------------------------------------------------------------|:--------------------|:----------------------------------------|
| [Queue](https://docs.kubemq.io/learn/message-patterns/queue)                      | target.queue        | [Usage](targets/queue/README.md)        |
| [Events](https://docs.kubemq.io/learn/message-patterns/pubsub#events)             | target.events       | [Usage](targets/events/README.md)       |
| [Events Store](https://docs.kubemq.io/learn/message-patterns/pubsub#events-store) | target.events-store | [Usage](targets/events-store/README.md) |
| [Command](https://docs.kubemq.io/learn/message-patterns/rpc#commands)             | target.command      | [Usage](targets/command/README.md)      |
| [Query](https://docs.kubemq.io/learn/message-patterns/rpc#queries)                | target.query        | [Usage](targets/query/README.md)        |

## Installation

### Kubernetes

1. Install KubeMQ Cluster

```bash
kubectl apply -f https://get.kubemq.io/deploy
```

2. Run KubeMQ Source deployment yaml

```bash
kubectl apply -f https://raw.githubusercontent.com/kubemq-hub/kubemq-sources/master/deploy-example.yaml
```

### Binary (Cross-platform)

Download the appropriate version for your platform from KubeMQ Targets Releases. Once downloaded, the binary can be run from anywhere.

Ideally, you should install it somewhere in your PATH for easy use. /usr/local/bin is the most probable location.

Running KubeMQ Targets

```bash
kubemq-targets --config config.yaml
```


## Configuration

### Structure

Config file structure:

```yaml

apiPort: 8080 # kubemq sources api and health end-point port
bindings:
  - name: clusters-sources # unique binding name
    properties: # Bindings properties such middleware configurations
      log_level: error
      retry_attempts: 3
      retry_delay_milliseconds: 1000
      retry_max_jitter_milliseconds: 100
      retry_delay_type: "back-off"
      rate_per_second: 100
    source:
      kind: source.http # source kind
      name: name-of-sources # source name 
      properties: # a set of key/value settings per each source kind
        .....
    target:
      kind: target.events # target kind
      name: name-of-target # targets name
      properties: # a set of key/value settings per each target kind
        - .....
```

### Properties

In bindings configuration, KubeMQ targets support properties setting for each pair of source and target bindings.

These properties contain middleware information settings as follows:

#### Logs Middleware

KubeMQ targets support level based logging to console according to as follows:

| Property  | Description       | Possible Values        |
|:----------|:------------------|:-----------------------|
| log_level | log level setting | "debug","info","error" |
|           |                   |  "" - indicate no logging on this bindings |

An example for only error level log to console:

```yaml
bindings:
  - name: sample-binding 
    properties: 
      log_level: error
    source:
    ......  
```

#### Retry Middleware

KubeMQ targets support Retries' target execution before reporting of error back to the source on failed execution.

Retry middleware settings values:


| Property                      | Description                                           | Possible Values                             |
|:------------------------------|:------------------------------------------------------|:--------------------------------------------|
| retry_attempts                | how many retries before giving up on target execution | default - 1, or any int number              |
| retry_delay_milliseconds      | how long to wait between retries in milliseconds      | default - 100ms or any int number           |
| retry_max_jitter_milliseconds | max delay jitter between retries                      | default - 100ms or any int number           |
| retry_delay_type              | type of retry delay                                   | "back-off" - delay increase on each attempt |
|                               |                                                       | "fixed" - fixed time delay                  |
|                               |                                                       | "random" - random time delay                |

An example for 3 retries with back-off strategy:

```yaml
bindings:
  - name: sample-binding 
    properties: 
      retry_attempts: 3
      retry_delay_milliseconds: 1000
      retry_max_jitter_milliseconds: 100
      retry_delay_type: "back-off"
    source:
    ......  
```

#### Rate Limiter Middleware

KubeMQ targets support a Rate Limiting of target executions.

Rate Limiter middleware settings values:


| Property        | Description                                    | Possible Values                |
|:----------------|:-----------------------------------------------|:-------------------------------|
| rate_per_second | how many executions per second will be allowed | 0 - no limitation              |
|                 |                                                | 1 - n integer times per second |

An example for 100 executions per second:

```yaml
bindings:
  - name: sample-binding 
    properties: 
      rate_per_second: 100
    source:
    ......  
```

### Source

Source section contains source configuration for Binding as follows:

| Property    | Description                                       | Possible Values                                               |
|:------------|:--------------------------------------------------|:--------------------------------------------------------------|
| name        | sources name (will show up in logs)               | string without white spaces                                   |
| kind        | source kind type                                  | source.queue                                                  |
|             |                                                   | source.query                                                  |
|             |                                                   | source.command                                                |
|             |                                                   | source.events                                                 |
|             |                                                   | source.events-store                                           |
| properties | an array of key/value setting for source connection| see above               |


### Target

Target section contains the target configuration for Binding as follows:

| Property    | Description                                       | Possible Values                                               |
|:------------|:--------------------------------------------------|:--------------------------------------------------------------|
| name        | targets name (will show up in logs)               | string without white spaces                                   |
| kind        | source kind type                                  | target.type-of-target                                                  |
| properties | an array of key/value set for target connection | see above              |





