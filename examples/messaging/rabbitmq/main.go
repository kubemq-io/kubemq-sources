package main

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"github.com/streadway/amqp"
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
		eventsCh, err := client.SubscribeToEvents(ctx, "events.messaging.rabbitmq", "", errCh)
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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(fmt.Errorf("error dialing rabbitmq, %w", err))
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(fmt.Errorf("error getting rabbitmq channel, %w", err))
	}
	msg := amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		DeliveryMode:    uint8(1),
		Priority:        uint8(0),
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "10000",
		Body:            []byte("rabbitmq data"),
	}
	err = channel.Publish("", "some-queue", false, false, msg)
	if err != nil {
		log.Fatal(fmt.Errorf("error publishing rabbitmq message, %w", err))
	}
	time.Sleep(3 * time.Second)
}
