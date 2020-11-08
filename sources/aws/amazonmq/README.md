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

| Properties Key                  | Required| Description                                             | Example                                                                |
|:--------------------------------|:--------|:--------------------------------------------------------|:-----------------------------------------------------------------------|
| host                            | yes     | AmazonMQ connection host (stomp+ssl endpoint)           | "localhost:1883" |
| dynamic_mapping                 | yes     | set if to map amazonmq Destination to kubemq channel    | "true"          |
| username                        | no      | set AmazonMQ username                                   | "username" |
| password                        | no      | set AmazonMQ password                                   | "password" |
| destination                     | yes     | set destination name                                    | "destination"         |


Example:

```yaml
    bindings:
    - name: amazonmq
      source:
        kind: aws.amazonmq
        properties:
          destination: some-queue
          host: localhost:61613
          password: admin
          username: admin
      target:
        kind: kubemq.events
        properties:
          address: localhost:50000
          auth_token: ""
          channel: event.aws.amazonmq
          client_id: test
      properties: {}

```
