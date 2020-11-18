module github.com/kubemq-hub/kubemq-sources

go 1.15

require (
	cloud.google.com/go/pubsub v1.6.2
	github.com/Azure/azure-event-hubs-go/v3 v3.3.2
	github.com/Azure/azure-service-bus-go v0.10.6
	github.com/Shopify/sarama v1.27.0
	github.com/aws/aws-sdk-go v1.34.31
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fortytw2/leaktest v1.3.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.0
	github.com/go-resty/resty/v2 v2.3.0 // indirect
	github.com/go-stomp/stomp v2.0.6+incompatible
	github.com/json-iterator/go v1.1.10
	github.com/kubemq-hub/builder v0.5.9
	github.com/kubemq-hub/ibmmq-sdk v0.3.5
	github.com/kubemq-io/kubemq-go v1.4.4
	github.com/labstack/echo/v4 v4.1.17
	github.com/nats-io/nuid v1.0.1
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.16.0
	google.golang.org/api v0.32.0
)

//replace github.com/kubemq-hub/builder => ../builder
