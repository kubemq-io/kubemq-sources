package main

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/go-stomp/stomp"
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
		eventsCh, err := client.SubscribeToEvents(ctx, "event.aws.amazonmq", "", errCh)
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

	netConn, err := tls.Dial("tcp", "localhost:61614", &tls.Config{})
	if err != nil {
		log.Fatal(err)
	}

	conn, err := stomp.Connect(netConn, stomp.ConnOpt.Login("admin", "admin"))
	if err != nil {
		log.Fatal(err)
	}

	conn.Send("my-dest", "text/plain", []byte("test-send-amazonmq"))
	time.Sleep(3 * time.Second)
}
