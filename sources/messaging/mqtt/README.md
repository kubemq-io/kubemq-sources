# Kubemq MQTT Source

Kubemq MQTT source provides a MQTT subscriber for processing MQTT's messages.

## Prerequisites
The following are required to run events source connector:

- kubemq cluster
- kubemq-sources deployment
- MQTT Broker


## Configuration

MQTT source connector configuration properties:

| Properties Key | Required | Description                    | Example          |
|:---------------|:---------|:-------------------------------|:-----------------|
| host           | yes      | mqtt connection host           | "localhost:1883" |
| topic          | yes      | set mqtt subscription topic    | "queue"          |
| dynamic_mapping          | yes      | set if to map mqtt topic to kubemq channel    | "true"          |
| username       | no       | set mqtt username              | "username"       |
| password       | no       | set mqtt password              | "password"       |
| client_id      | no       | mqtt connection string address | "client_id"      |
| qos            | no       | set mqtt subscription QoS      | "0"              |


Example:

```yaml
bindings:
  - name: mqtt-kubemq-event
    source:
      kind: messaging.mqtt
      name: mqtt-source
      properties:
        host: "localhost:1883"
        dynamic_map: "true"
        topic: "queue"
        username: "username"
        password: "password"
        client_id: "client_id"
        qos: "0"
    target:
      kind: kubemq.events
      name: target-kubemq-events
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-http-connector"
        channel: "events.mqtt"
    properties:
      log_level: "info"
```
