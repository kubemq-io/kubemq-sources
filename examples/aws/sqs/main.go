package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	dat, err := ioutil.ReadFile("./credentials/aws/awsKey.txt")
	if err != nil {
		log.Fatal(err)
	}
	awsKey := string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/awsSecretKey.txt")
	if err != nil {
		log.Fatal(err)
	}
	awsSecretKey := string(dat)
	dat, err = ioutil.ReadFile("./credentials/aws/region.txt")
	if err != nil {
		log.Fatal(err)
	}
	region := string(dat)

	dat, err = ioutil.ReadFile("./credentials/aws/sqs/queue.txt")
	if err != nil {
		log.Fatal(err)
	}
	queue := string(dat)
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
		eventsCh, err := client.SubscribeToEvents(ctx, "event.aws.sqs", "", errCh)
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
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(awsKey, awsSecretKey, ""),
	})
	if err != nil {
		log.Fatal(err)
	}

	svc := sqs.New(sess)

	m := &sqs.SendMessageInput{}
	m.SetQueueUrl(queue)
	m.SetMessageBody("my test body!")
	_, err = svc.SendMessage(m)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(10 * time.Second)
	m2 := &sqs.SendMessageInput{}
	m2.SetQueueUrl(queue)
	m2.SetMessageBody("my second body!")
	_, err = svc.SendMessage(m2)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(10 * time.Second)
}
