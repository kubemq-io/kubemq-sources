# Kubemq IBM-MQ Source Connector

Kubemq IBM-MQ Source connector allows services using kubemq server to access IBM-MQ messaging services.

## Prerequisites
The following are required to run the IBM-MQ Source connector:

- kubemq cluster
- IBM-MQ server
- kubemq-Sources deployment

## Configuration

IBM-MQ Source connector configuration properties:

| Properties Key       | Required | Description                                          | Example                                                                |
|:---------------------|:---------|:-----------------------------------------------------|:-----------------------------------------------------------------------|
| queue_manager_name   | yes      | queue manager name (QMName)                          | "QM1"                                            |
| host_name            | yes      | set the host address or name                         | "localhost"                                      |
| channel_name         | yes      | ibm mq ChannelName to search the queue at            | "DEV.APP.CHANNEL"                                |
| username             | yes      | Username to use to login to ibm-mq                   | "my_user"                                        |
| queue_name           | yes      | the queue name to search under the channel           | "my_queue"                                       |
| certificate_label    | no       | unique identifier representing a digital certificate | "certificate label"                              |
| transport_type       | no       | ibmmq transport_type                                 | "0"(TransportType_CLIENT),"1"(TransportType_BINDINGS") |
| tls_client_auth      | no       | tls client auth type                                 | "NONE","REQUIRED"                                |
| port_number          | no       | ibmmq - server port number (default 1414)            | "1414"                                           |
| password             | no       | set ibmmq user password                              | "password"                                       |
| key_repository       | no       | the key_repository a certificate store               | "as123aq"                                        |


Example:

```yaml
bindings:
  - name: kubemq-query-IBM-MQ
    source:
      kind: kubemq.query
      name: kubemq-query
      properties:
        address: "kubemq-cluster:50000"
        client_id: "kubemq-query-IBM-MQ-connector"
        auth_token: ""
        channel: "query.messaging.ibmmq"
        group:   ""
        auto_reconnect: "true"
        reconnect_interval_seconds: "1"
        max_reconnects: "0"
    Source:
      kind: messaging.ibmmq
      name: messaging-ibmmq
      properties:
        queue_manager_name: "QM1"
        host_name: "localhost"
        channel_name: "DEV.APP.SVRCONN"
        username: "app"
        queue_name: "admin"
        password: "passw0rd"
        certificate_label: "NONE"
```
