# Kubemq servicebus source Connector

Kubemq azure-servicebus source connector allows services using kubemq server to access google servicebus server.

## Prerequisites
The following required to run the azure-servicebus source connector:

- kubemq cluster
- azure-servicebus set up
- kubemq-source deployment

## Configuration

event-hubs source connector configuration properties:

| Properties Key                  | Required | Description                                                            | Example                                                                |
|:--------------------------------|:---------|:-----------------------------------------------------------------------|:-----------------------------------------------------------------------|
| end_point                       | yes      | servicebus target endpoint                                             | "sb://my_account.net" |
| shared_access_key_name          | yes      | servicebus access key name                                             | "keyname" |
| shared_access_key               | yes      | servicebus shared access key name                                      | "213ase123" |
| queue                           | yes      | servicebus queue name                                                  | "0" |




Example:

```yaml
bindings:
  - name: kubemq-query-azure-servicebus
    source:
      kind: query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-azure-servicebus-connector"
        auth_token: ""
        channel: "event.azure.servicebus"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: azure.servicebus
      name: azure-servicebus
      properties:
          end_point: "sb://my_account.net"
          shared_access_key_name: "keyname"
          shared_access_key: "213ase123"
          queue: "test"  

```
