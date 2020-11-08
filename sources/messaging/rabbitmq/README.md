# Kubemq RabbitMQ Source

Kubemq RabbitMQ source provides a RabbitMQ subscriber for processing RabbitMQ's messages.

## Prerequisites
The following are required to run events source connector:

- kubemq cluster
- kubemq-sources deployment
- RabbitMQ Server


## Configuration

RabbitMQ source connector configuration properties:

| Properties Key   | Required | Description                         | Example                                    |
|:-----------------|:---------|:------------------------------------|:-------------------------------------------|
| url              | yes      | rabbitmq connection string address  | "amqp://guest:guest@localhost:5672/" |
| queue            | yes      | set subscription queue              | "queue"                                    |
| dynamic_mapping          | yes      | set if to map rabbit topic to kubemq channel    | "true"          |
| consumer         | yes      | set subscription consumer tag       | "consumer"                                 |
| requeue_on_error | bool     | set requeue on error property       | "false"                                    |
| auto_ack         | bool     | set auto_ack property               | "false"                                    |
| exclusive        | bool     | set exclusive property              | "false"                                    |


Example:

```yaml
bindings:
- name: rabbitmq
  source:
    kind: messaging.rabbitmq
    properties:
      auto_ack: "false"
      consumer: "1"
      exclusive: "false"
      dynamic_mapping: "true"
      queue: some-queue
      requeue_on_error: "false"
      url: amqp://guest:guest@localhost:5672/
  target:
    kind: kubemq.events
    properties:
      address: localhost:50000
      auth_token: ""
      channel: events.messaging.rabbitmq
      client_id: rabbitmq
  properties: {}

```
