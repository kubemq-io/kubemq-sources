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
| project_id     | yes      | gcp project_id                             | "<googleurl>/myproject"    |
| credentials    | yes      | gcp credentials files                      | "<google json credentials" |
| subscriber_id  | yes      | gcp pubsub queue subscriber id             | "string"          |


Example:

```yaml
bindings:
  - name: kubemq-query-gcp-pubsub
    source:
      kind: query
      name: kubemq-query
      properties:
        host: "localhost"
        port: "50000"
        client_id: "kubemq-query-gcp-pubsub-connector"
        auth_token: ""
        channel: "query.gcp.pubsub"
        group:   ""
        concurrency: "1"
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: gcp.pubsub
      name: source-gcp-pubsub
      properties:
        project_id: "projectID"
        subscriber_id:    "mysubscriberID"
        credentials: 'json'

```
