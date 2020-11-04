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
| consumer         | yes      | set subscription consumer tag       | "consumer"                                 |
| requeue_on_error | bool     | set requeue on error property       | "false"                                    |
| auto_ack         | bool     | set auto_ack property               | "false"                                    |
| exclusive        | bool     | set exclusive property              | "false"                                    |


Example:

```yaml
bindings:
  - name: rabbitmq-kubemq-event
    source:
      kind: messaging.rabbitmq
      name: rabbitmq-source
      properties:
        "url": "amqp://guest:guest@localhost:5672/"
        "queue": "some-queue"
        "consumer": "kubemq"
        "requeue_on_error": "false"
        "auto_ack": "false"
        "exclusive": "false"
    target:
      kind: kubemq.events
      name: target-kubemq-events
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-http-connector"
        channel: "events.rabbitmq"
    properties:
      log_level: "info"
```
