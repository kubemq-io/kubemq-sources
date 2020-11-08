# KubeMQ Sources Query Target

KubeMQ Sources Query target provides an RPC query sender for processing sources requests.

## Prerequisites
The following are required to run the query target connector:

- kubemq cluster
- kubemq-sources deployment


## Configuration

Query target connector configuration properties:

| Properties Key  | Required | Description                                        | Example                                              |
|:----------------|:---------|:---------------------------------------------------|:-----------------------------------------------------|
| address         | yes      | kubemq server address (gRPC interface)             | kubemq-cluster-grpc.kubemq.svc.cluster.local:50000 |
| client_id       | no       | set client id                                      | "client_id"                                          |
| auth_token      | no       | set authentication token                           | JWT token                                            |
| channel | no       | set send request default channel               |          "queries"                                            |
| dynamic_mapping | no       | set dynamic channel mapping per source               |          "false"                                            |
| timeout_seconds | no       | sets query request default timeout (600 seconds) |          "10"                                            |


Example:

```yaml
bindings:
  - name:  query-binding 
    properties: 
      log_level: error
      retry_attempts: 3
      retry_delay_milliseconds: 1000
      retry_max_jitter_milliseconds: 100
      retry_delay_type: "back-off"
      rate_per_second: 100
    source:
      kind: http
      name: http-get-source
      properties:
        "methods": "get"
        "path": "/query"
    target:
      kind: kubemq.query # Sources kind
      name: query-target 
      properties: 
        address: "kubemq-cluster-grpc.kubemq.svc.cluster.local:50000"
        client_id: "cluster-a-query-connection"
        auth_token: ""
        channel: "queries"
        dynamic_mapping: "false"
        timeout_seconds: "10"
```

