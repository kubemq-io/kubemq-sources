package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
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
		eventsCh, err := client.SubscribeToEvents(ctx, "events.>", "", errCh)
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

	time.Sleep(time.Second)
	req, err := http.NewRequest("POST", "http://localhost:8090/post", bytes.NewBufferString("post-data"))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	time.Sleep(3 * time.Second)
}
