# Kubemq amazonMQ Source Connector

Kubemq AmazonMQ Source connector allows services using kubemq server to access AmazonMQ messaging services.

## Prerequisites
The following are required to run the AmazonMQ Source connector:

- kubemq cluster
- AmazonMQ server - with access 
- kubemq-Sources deployment


- Please note the connector uses connection with stomp+ssl, when finishing handling messages need to call Stop().

## Configuration

AmazonMQ Source connector configuration properties:

| Properties Key                  | Required | Description                                 | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------|:-----------------------------------------------------------------------|
| host                            | yes     | AmazonMQ connection host (stomp+ssl endpoint)| "localhost:1883" |
| username                        | no      | set AmazonMQ username                        | "username" |
| password                        | no      | set AmazonMQ password                        | "password" |
| destination                     | yes     | set destination name                         | "destination"         |
| subTimeout                      | no      | timeout for sub expired                      | "5"(default 5)         |


Example:

```yaml
bindings:
  - name: kubemq-query-amazonmq
    source:
      kind: query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-amazonmq-connector"
        auth_token: ""
        channel: "query.amazonmq"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: aws.amazonmq
      name: source-aws-amazonmq
      properties:
        host: "localhost:61613"
        username: "admin"
        password: "admin"
        destination: "my-queue"
        subTimeout: "5"
```
