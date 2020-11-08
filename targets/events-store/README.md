# KubeMQ Sources Events-Store Target

KubeMQ Sources Events-Store target provides an events-store sender for processing sources requests.

## Prerequisites
The following are required to run the events-store target connector:

- kubemq cluster
- kubemq-sources deployment


## Configuration

Events-Store target connector configuration properties:

| Properties Key  | Required | Description                                        | Example                                              |
|:----------------|:---------|:---------------------------------------------------|:-----------------------------------------------------|
| address         | yes      | kubemq server address (gRPC interface)             | kubemq-cluster-grpc.kubemq.svc.cluster.local:50000 |
| client_id       | no       | set client id                                      | "client_id"                                          |
| auth_token      | no       | set authentication token                           | JWT token                                            |
| channel | no       | set send request default channel               |          "events-store"                                            |
| dynamic_mapping | no       | set dynamic channel mapping per source               |          "false"                                            |

Example:

```yaml
bindings:
  - name:  events-store-binding 
    properties: 
      log_level: error
      retry_attempts: 3
      retry_delay_milliseconds: 1000
      retry_max_jitter_milliseconds: 100
      retry_delay_type: "back-off"
      rate_per_second: 100
    source:
      kind: http
      name: http-post-source
      properties:
        "methods": "post"
        "path": "/events-store"
    target:
      kind: kubemq.events-store # Sources kind
      name: events-store-target 
      properties: 
        address: "kubemq-cluster-grpc.kubemq.svc.cluster.local:50000"
        client_id: "cluster-a-events-store-connection"
        auth_token: ""
        channel: "events-store"
        dynamic_mapping: "false"
        timeout_seconds: "10"
```

