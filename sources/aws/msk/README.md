# Kubemq msk Source Connector

Kubemq msk source connector allows services using kubemq server to access msk server. TODO

## Prerequisites
The following are required to run the msk target connector:

- kubemq cluster
- active msk cluster
- kubemq-target-connectors deployment

## Configuration

msk source connector configuration properties:

| Properties Key | Required | Description                                | Example          |
|:---------------|:---------|:-------------------------------------------|:-----------------|
| brokers        | yes      | msk brokers connection, comma separated    | "localhost:9092" |
| topics         | yes      | msk stored topic, comma separated          | "TestTopic"      |
| dynamic_mapping| yes      | set if to map msk topic to kubemq channel  | "true"          |
| consumer_group | yes      | msk consumer group name                    | "Group1          |
| sasl_username  | no       | SASL based authentication with broker      | "user"           |
| sasl_password  | no       | SASL based authentication with broker      | "pass"           |

Example:

```yaml
bindings:
  - name: kubemq-store-msk
    source:
        kind: kubemq.event-store
        name: kubemq-query
        properties:
            host: "localhost"
            port: "50000"
            client_id: "kubemq-query-aws-connector"
            auth_token: ""
            channel: "aws.msk"
            group:   ""
            auto_reconnect: "true"
            reconnect_interval_seconds: "1"
            max_reconnects: "0"
    target:
      kind: aws.msk
      name: source-aws-msk
      properties:
     	brokers: "localhost:9092,localhost:9093",
	    topic: "TestTopic",
	    consumer_group: "cg"
```
