package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/azure/eventhubs/connectionString.txt")
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
		eventsCh, err := client.SubscribeToEvents(ctx, "event.azure.eventhubs", "", errCh)
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
	eventHubClient, err := eventhub.NewHubFromConnectionString(connectionString)
	event := &eventhub.Event{
		Data: []byte("new test"),
	}
	err = eventHubClient.Send(ctx, event)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(30 * time.Second)
}
