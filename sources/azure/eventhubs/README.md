# Kubemq eventhubs source Connector

Kubemq azure-eventhubs source connector allows services using kubemq server to access google eventhubs server.

## Prerequisites
The following required to run the azure-eventhubs source connector:

- kubemq cluster
- azure-eventhubs set up
- kubemq-source deployment

## Configuration

event-hubs source connector configuration properties:

| Properties Key                  | Required | Description                                                            | Example                                                                |
|:--------------------------------|:---------|:-----------------------------------------------------------------------|:-----------------------------------------------------------------------|
| end_point                       | yes      | event hubs target endpoint                                             | "sb://my_account.net" |
| shared_access_key_name          | yes      | event hubs access key name                                             | "keyname" |
| shared_access_key               | yes      | event hubs shared access key name                                      | "213ase123" |
| entity_path                     | yes      | event hubs path entity to subscribe                                    | "mypath" |
| partition_id                    | yes      | event hubs partition_id to listen                                      | "0" |
| receive_type                    | no       | event hubs path entity to send                                         | "latest_offset","from_timestamp","with_consumer_group","with_epoch","with_prefetch_count","with_starting_offset" Default(with_starting_offset) |
| partition_id                    | yes      | event hubs partition_id to listen                                      | "0" |
| time_stamp                      | no       | timestamp to read events from must supply with from_timestamp(RFC3339) | "0" |
| with_consumer_group             | no       | consumer group to assign the listener must supply with with_consumer_group(RFC3339)  | "0" |
| with_epoch                      | no       | timestamp to read from w must supply with with_epoch(RFC3339)                             | "0" |
| with_prefetch_count             | no       | event hubs partition_id to listen                                    | "0" |
| with_starting_offset            | no       | event hubs partition_id to listen                                    | "0" |




Example:

```yaml
bindings:
  - name: kubemq-query-azure-eventhubs
    source:
      kind: query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-azure-eventhubs-connector"
        auth_token: ""
        channel: "query.azure.eventhubs"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    target:
      kind: azure.eventhubs
      name: azure-eventhubs
      properties:
          end_point: "sb://my_account.net"
          shared_access_key_name: "keyname"
          shared_access_key: "213ase123"
          entity_path: "mypath"
          partition_id: "0"  

```
