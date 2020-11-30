# Kubemq nats source Connector

Kubemq nats source connector allows services using kubemq server to access nats messaging services.

## Prerequisites
The following are required to run the nats source connector:

- kubemq cluster
- nats server
- kubemq-sources deployment

## Configuration

nats source connector configuration properties:

| Properties Key                  | Required | Description                                             | Example                                                                |
|:--------------------------------|:---------|:--------------------------------------------------------|:-----------------------------------------------------------------------|
| url                             | yes      | nats connection host                                    | "localhost:1883" |
| subject                         | yes      | set subject name                                        | any string |
| dynamic_mapping                 | yes      | set if to map nats Destination to kubemq channel        | "true"          |
| username                        | no       | set nats username                                       | "username" |
| password                        | no       | set nats password                                       | "password" |
| token                           | no       | set nats token                                          | "my_token" |
| tls                             | no       | set if tls is needed                                    | "false","true" |
| cert_file                       | no       | tls certificate file in string format                   | "my_file" |
| cert_key                        | no       | tls certificate key in string format                    | "my_key"  |
| timeout                         | no       | connection timeout in seconds                           | "130"  |


Example:

```yaml
bindings:
  - name: nats
    source:
      kind: messaging.nats
      properties:
        cert_file: |-
          -----BEGIN CERTIFICATE-----
          mycert
          -----END CERTIFICATE-----
        cert_key: |-
          -----BEGIN PRIVATE KEY-----
          mykey
          -----END PRIVATE KEY-----
        dynamic_mapping: "false"
        subject: foo
        url: nats://localhost:4222
    target:
      kind: kubemq.events
      properties:
        address: localhost:50000
        channel: event.messaging.nats
        dynamic_mapping: "false"
    properties: {}

```