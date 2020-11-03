package main

import (
	"context"
	kafka "github.com/Shopify/sarama"
	"github.com/kubemq-io/kubemq-go"
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
		eventsCh, err := client.SubscribeToEvents(ctx, "event.messaging.msk", "", errCh)
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

	kc := kafka.NewConfig()
	kc.Version = kafka.V2_0_0_0
	kc.Producer.RequiredAcks = kafka.WaitForAll
	kc.Producer.Retry.Max = 5
	kc.Producer.Return.Successes = true
	var brokers []string
	brokers = append(brokers, "localhost:9092")
	producer, err := kafka.NewSyncProducer(brokers, kc)
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = producer.SendMessage(&kafka.ProducerMessage{
		Key:   kafka.ByteEncoder("test"),
		Value: kafka.ByteEncoder("test_value"),
		Topic: "TestTopicA",
	})
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
}
