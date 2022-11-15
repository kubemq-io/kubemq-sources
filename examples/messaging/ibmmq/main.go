package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/kubemq-hub/ibmmq-sdk/mq-golang-jms20/mqjms"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
)

type testStructure struct {
	applicationChannelName string
	hostname               string
	listenerPort           string
	queueManagerName       string
	apiKey                 string
	mqUsername             string
	password               string
	QueueName              string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./credentials/ibm/mq/connectionInfo/applicationChannelName.txt")
	if err != nil {
		return nil, err
	}
	t.applicationChannelName = string(dat)
	dat, err = ioutil.ReadFile("./credentials/ibm/mq/connectionInfo/hostname.txt")
	if err != nil {
		return nil, err
	}
	t.hostname = string(dat)
	dat, err = ioutil.ReadFile("./credentials/ibm/mq/connectionInfo/listenerPort.txt")
	if err != nil {
		return nil, err
	}
	t.listenerPort = string(dat)
	dat, err = ioutil.ReadFile("./credentials/ibm/mq/connectionInfo/queueManagerName.txt")
	if err != nil {
		return nil, err
	}
	t.queueManagerName = string(dat)

	dat, err = ioutil.ReadFile("./credentials/ibm/mq/applicationApiKey/apiKey.txt")
	if err != nil {
		return nil, err
	}
	t.apiKey = string(dat)
	dat, err = ioutil.ReadFile("./credentials/ibm/mq/applicationApiKey/mqUsername.txt")
	if err != nil {
		return nil, err
	}
	t.mqUsername = string(dat)
	dat, err = ioutil.ReadFile("./credentials/ibm/mq/applicationApiKey/mqPassword.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)
	t.QueueName = "DEV.QUEUE.1"
	return t, nil
}

func main() {
	dat, err := getTestStructure()
	if err != nil {
		log.Fatal(err)
	}
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
		eventsCh, err := client.SubscribeToEvents(ctx, "events.messaging.ibmmq", "", errCh)
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

	cf := mqjms.ConnectionFactoryImpl{
		QMName:        dat.queueManagerName,
		Hostname:      dat.hostname,
		PortNumber:    1414,
		ChannelName:   dat.applicationChannelName,
		UserName:      dat.mqUsername,
		TransportType: 0,
		TLSClientAuth: "NONE",
		KeyRepository: "",
		Password:      dat.password,

		CertificateLabel: "",
	}

	jmsContext, jmsErr := cf.CreateContext()
	if jmsErr != nil {
		log.Fatal(jmsErr.GetReason())
	}
	queue := jmsContext.CreateQueue(dat.QueueName)
	producer := jmsContext.CreateProducer()
	producer.SendString(queue, "test")
	time.Sleep(3 * time.Second)
}
