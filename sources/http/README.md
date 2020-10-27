# Kubemq HTTP Source

Kubemq HTTP source provides a REST Api endpoints for processing source REST Api request.

## Prerequisites
The following are required to run events source connector:

- kubemq cluster
- kubemq-sources deployment


## Configuration

HTTP source connector configuration properties:

| Properties Key             | Required | Description                           | Example            |
|:---------------------------|:---------|:--------------------------------------|:-------------------|
| methods                    | yes      | List for REST API methods to process| "GET,POST,PUT,PATCH,DELETE"|
| path                  | yes      | set http server endpoint path                      | "/path"        |


Example:

```yaml
bindings:
  - name: http-post-kubemq-event
    source:
      kind: http
      name: http-post
      properties:
        "methods": "post"
        "path": "/post"
    target:
      kind: kubemq.events
      name: target-kubemq-events
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-http-connector"
        auth_token: ""
        channel: "events.post"
    properties:
      log_level: "info"
```
