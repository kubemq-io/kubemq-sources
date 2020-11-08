# KubeMQ Sources Queue Target

KubeMQ Sources Queue target provides a queue sender for processing sources requests.

## Prerequisites
The following are required to run the queue target connector:

- kubemq cluster
- kubemq-sources deployment


## Configuration

Queue target connector configuration properties:

| Properties Key  | Required | Description                                        | Example                                              |
|:----------------|:---------|:---------------------------------------------------|:-----------------------------------------------------|
| address         | yes      | kubemq server address (gRPC interface)             | kubemq-cluster-grpc.kubemq.svc.cluster.local:50000 |
| client_id       | no       | set client id                                      | "client_id"                                          |
| auth_token      | no       | set authentication token                           | JWT token                                            |
| channel | no       | set send request default channel               |          "queue"                                            |
| dynamic_mapping | no       | set dynamic channel mapping per source               |          "false"                                            |



Example:

```yaml
bindings:
  - name:  queue-binding 
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
        "path": "/queue"
    target:
      kind: kubemq.queue # Sources kind
      name: queue-target 
      properties: 
        address: "kubemq-cluster-grpc.kubemq.svc.cluster.local:50000"
        client_id: "cluster-a-queue-connection"
        auth_token: ""
        channel: "queue"
        dynamic_mapping: "false"
        timeout_seconds: "10"
```

