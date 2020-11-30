package main

import (
	"context"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
	"log"
	"time"
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
		eventsCh, err := client.SubscribeToEvents(ctx, "event.messaging.nats", "", errCh)
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
				log.Printf("client: Channel: %s, Data: %s", event.Channel, string(event.Body))

			case <-ctx.Done():
				return
			}
		}

	}()

	nc, err := nats.Connect("nats://localhost:4222", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subj, msg := "foo", []byte("test")

	err = nc.Publish(subj, msg)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
}
