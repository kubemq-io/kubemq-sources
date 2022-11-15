# Kubemq Kafka Source Connector

Kubemq kafka source connector allows services using kubemq server to access redis server. TODO

## Prerequisites
The following are required to run the redis target connector:

- kubemq cluster
- kafka TODO version
- kubemq-sources deployment

## Configuration

Kafka source connector configuration properties:

| Properties Key     | Required | Description                                 | Example                                                       |
|:-------------------|:---------|:--------------------------------------------|:--------------------------------------------------------------|
| brokers            | yes      | kafka brokers connection, comma separated   | "localhost:9092"                                              |
| topics             | yes      | kafka stored topic, comma separated         | "TestTopic"                                                   |
| dynamic_mapping    | yes      | set if to map kafka topic to kubemq channel | "true"                                                        |
| consumer_group     | yes      | kafka consumer group name                   | "Group1                                                       |
| saslUsername       | no       | SASL based authentication with broker       | "user"                                                        |
| saslPassword       | no       | SASL based authentication with broker       | "pass"                                                        |
| saslMechanism      | no       | SASL Mechanism                              | SCRAM-SHA-256, SCRAM-SHA-512, plain, 0Auth bearer, or GSS-API |
| securityProtocol   | no       | Set connection security protocol            | plaintext, SASL-plaintext, SASL-SSL, SSL                      |
| ca_cert            | no       | SSL CA certificate                          | pem certificate value                                         |
| client_certificate | no       | SSL Client certificate (mMTL)               | pem certificate value                                         |
| client_key         | no       | SSL Client Key (mTLS)                       | pem key value                                                 |


Example:

```yaml
bindings:
  - name: kafka
    source:
      kind: messaging.kafka
      properties:
        brokers: localhost:9092
        consumer_group: test_client
        topics: TestTopicA
    target:
      kind: kubemq.events
      properties:
        address: localhost:50000
        auth_token: ""
        channel: event.messaging.kafka
        client_id: test
    properties: {}

```

## Usage

