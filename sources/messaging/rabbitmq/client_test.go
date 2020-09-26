package rabbitmq

//
//import (
//	"context"
//	"github.com/kubemq-hub/kubemq-sources/config"
//	"github.com/kubemq-hub/kubemq-sources/middleware"
//	"github.com/kubemq-hub/kubemq-sources/targets"
//	"github.com/kubemq-hub/kubemq-sources/targets/null"
//	"github.com/kubemq-hub/kubemq-sources/types"
//	"github.com/nats-io/nuid"
//	"github.com/streadway/amqp"
//	"github.com/stretchr/testify/require"
//	"testing"
//	"time"
//)
//
//var rabbitmqUrl = "amqp://rabbitmq:rabbitmq@localhost:5672/"
//
//func setupClient(ctx context.Context, queue string, target middleware.Middleware) (*Client, error) {
//	c := New()
//	err := c.Init(ctx, config.Spec{
//		Name: "rabbitmq",
//		Kind: "",
//		Properties: map[string]string{
//			"url":              rabbitmqUrl,
//			"queue":            queue,
//			"consumer":         nuid.Next(),
//			"requeue_on_error": "false",
//			"auto_ack":         "false",
//			"exclusive":        "false",
//		},
//	})
//	if err != nil {
//		return nil, err
//	}
//	err = c.Start(ctx, target)
//	if err != nil {
//		return nil, err
//	}
//	time.Sleep(time.Second)
//	return c, nil
//}
//
//func sendMessage(queue string, data []byte) error {
//	conn, err := amqp.Dial(rabbitmqUrl)
//	if err != nil {
//		return err
//	}
//	channel, err := conn.Channel()
//	if err != nil {
//		return err
//	}
//	err = channel.Publish("", queue, false, false, amqp.Publishing{
//		Headers:         amqp.Table{},
//		ContentType:     "text/plain",
//		ContentEncoding: "",
//		DeliveryMode:    1,
//		Priority:        0,
//		CorrelationId:   "",
//		ReplyTo:         "",
//		Expiration:      "",
//		MessageId:       "",
//		Timestamp:       time.Time{},
//		Type:            "",
//		UserId:          "",
//		AppId:           "",
//		Body:            data,
//	})
//	return err
//}
//
//func TestClient_Start(t *testing.T) {
//	tests := []struct {
//		name    string
//		target  -sources.Target
//		req     *types.Request
//		queue   string
//		wantErr bool
//	}{
//		{
//			name: "request",
//			target: &null.Client{
//				Delay:         0,
//				DoError:       nil,
//				ResponseError: nil,
//			},
//			req:     types.NewRequest().SetData([]byte("some-data")),
//			queue:   "some-queue",
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			c, err := setupClient(ctx, tt.queue, tt.target)
//			require.NoError(t, err)
//			defer func() {
//				_ = c.Stop()
//			}()
//
//			err = sendMessage(tt.queue, tt.req.Data)
//			if tt.wantErr {
//				require.Error(t, err)
//				return
//			}
//			require.NoError(t, err)
//		})
//	}
//}
//
//func TestClient_Init(t *testing.T) {
//
//	tests := []struct {
//		name    string
//		cfg     config.Spec
//		wantErr bool
//	}{
//		{
//			name: "init",
//			cfg: config.Spec{
//				Name: "rabbitmq",
//				Kind: "",
//				Properties: map[string]string{
//					"url":              "amqp://rabbitmq:rabbitmq@localhost:5672/",
//					"queue":            "some-queue",
//					"consumer":         nuid.Next(),
//					"requeue_on_error": "false",
//					"auto_ack":         "false",
//					"exclusive":        "false",
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "init - bad url",
//			cfg: config.Spec{
//				Name: "rabbitmq-target",
//				Kind: "",
//				Properties: map[string]string{
//					"url":              "amqp://rabbitmq:rabbitmq@localhost:6000/",
//					"consumer":         nuid.Next(),
//					"queue":            "some-queue",
//					"requeue_on_error": "false",
//					"auto_ack":         "false",
//					"exclusive":        "false",
//				},
//			},
//			wantErr: true,
//		},
//		{
//			name: "bad init - no  url",
//			cfg: config.Spec{
//				Name: "rabbitmq",
//				Kind: "",
//				Properties: map[string]string{
//					"queue":            "some-queue",
//					"consumer":         nuid.Next(),
//					"requeue_on_error": "false",
//					"auto_ack":         "false",
//					"exclusive":        "false",
//				},
//			},
//			wantErr: true,
//		},
//		{
//			name: "bad init - no queue",
//			cfg: config.Spec{
//				Name: "rabbitmq",
//				Kind: "",
//				Properties: map[string]string{
//					"url":              "amqp://rabbitmq:rabbitmq@localhost:5672/",
//					"consumer":         nuid.Next(),
//					"requeue_on_error": "false",
//					"auto_ack":         "false",
//					"exclusive":        "false",
//				},
//			},
//			wantErr: true,
//		},
//		{
//			name: "init - no consumer",
//			cfg: config.Spec{
//				Name: "rabbitmq-target",
//				Kind: "",
//				Properties: map[string]string{
//					"url":              "amqp://rabbitmq:rabbitmq@localhost:5432/",
//					"queue":            "some-queue",
//					"requeue_on_error": "false",
//					"auto_ack":         "false",
//					"exclusive":        "false",
//				},
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//			defer cancel()
//			c := New()
//			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
//				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
//			}
//
//		})
//	}
//}
