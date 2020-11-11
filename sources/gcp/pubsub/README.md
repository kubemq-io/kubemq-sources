# Kubemq pubsub source Connector

Kubemq gcp-pubsub source connector allows services using kubemq server to access google pubsub server.

## Prerequisites
The following required to run the gcp-pubsub source connector:

- kubemq cluster
- gcp-pubsub set up
- kubemq-source deployment

## Configuration

pubsub source connector configuration properties:

| Properties Key | Required | Description                                | Example                    |
|:---------------|:---------|:-------------------------------------------|:---------------------------|
| project_id     | yes      | gcp project_id                             | "googleurl/myproject"    |
| credentials    | yes      | gcp credentials files                      | "google json credentials" |
| subscriber_id  | yes      | gcp pubsub queue subscriber id             | "string"          |


Example:

```yaml
bindings:
  - name: pubsub
    source:
      kind: gcp.pubsub
      properties:
        credentials: |-
          {
          }
        project_id: kubemq-tests-294511
        subscriber_id: test
    target:
      kind: kubemq.events
      properties:
        address: localhost:50000
        auth_token: ""
        channel: event.gcp.pubsub
        client_id: pubsub
    properties: {}


```
