package main

import (
	"context"
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/azure/servicebus/connectionString.txt")
	if err != nil {
		log.Fatal(err)
	}
	connectionString := string(dat)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithCheckConnection(true),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		errCh := make(chan error)
		eventsCh, err := client.SubscribeToEvents(ctx, "event.azure.servicebus", "", errCh)
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case err := <-errCh:
				log.Fatal(err)
				return
			case event, more := <-eventsCh:
				if !more {
					log.Println("client done")
					return
				}
				log.Printf("client:%s Channel: %s, Data: %s", event.ClientId, event.Channel, string(event.Body))

			case <-ctx.Done():
				return
			}
		}

	}()
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	queueName := "myqueue"
	svcBus, err := ns.NewQueue(queueName)
	if err != nil {
		log.Fatal(err)
	}
	err = svcBus.Send(ctx, servicebus.NewMessage([]byte("my_new_message")))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(30 * time.Second)
}
