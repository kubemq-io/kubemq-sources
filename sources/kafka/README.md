# Kubemq Kafka Source Connector

Kubemq kafka source connector allows services using kubemq server to access redis server. TODO

## Prerequisites
The following are required to run the redis target connector:

- kubemq cluster
- kafka TODO version
- kubemq-target-connectors deployment

## Configuration

Kafka source connector configuration properties:

| Properties Key | Required | Description                                | Example          |
|:---------------|:---------|:-------------------------------------------|:-----------------|
| brokers        | yes      | kafka brokers connection, comma separated  | "localhost:9092" |
| topics         | yes      | kafka stored topic, comma separated        | "TestTopic"      |
| consumerGroup  | yes      | kafka consumer group name                  | "Group1          |
| saslUsername   | no       | SASL based authentication with broker      | "user"           |
| saslPassword   | no       | SASL based authentication with broker      | "pass"           |

Example:

```yaml
bindings:
  - name: kubemq-store-kafka
    source:
        kind: source.kubemq.event-store
        name: kubemq-query
        properties:
            host: "localhost"
            port: "50000"
            client_id: "kubemq-query-redis-connector"
            auth_token: ""
            channel: "store.kafka"
            group:   ""
            auto_reconnect: "true"
            reconnect_interval_seconds: "1"
            max_reconnects: "0"
    target:
      kind: target.messaging.kafka
      name: kafka-stream
      properties:
     	brokers: "localhost:9092,localhost:9093",
		topic: "TestTopic",
		consumerGroup: "cg"
```

## Usage

