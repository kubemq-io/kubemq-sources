package main

import (
	"context"
	"log"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
)

func main() {
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
		eventsCh, err := client.SubscribeToEvents(ctx, "event.gcp.pubsub", "", errCh)
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
	time.Sleep(10 * time.Second)
}
