package mqtt
//
//import (
//	"context"
//	"fmt"
//	mqtt "github.com/eclipse/paho.mqtt.golang"
//	"github.com/kubemq-hub/kubemq-sources/config"
//	"github.com/kubemq-hub/kubemq-sources/middleware"
//	"github.com/kubemq-hub/kubemq-sources/targets"
//	"github.com/kubemq-hub/kubemq-sources/targets/null"
//	"github.com/kubemq-hub/kubemq-sources/types"
//	"github.com/nats-io/nuid"
//	"github.com/stretchr/testify/require"
//	"testing"
//	"time"
//)
//
//func setupClient(ctx context.Context, queue string, target middleware.Middleware) (*Client, error) {
//	c := New()
//	err := c.Init(ctx, config.Spec{
//		Name: "mqtt",
//		Kind: "",
//		Properties: map[string]string{
//			"host":     "localhost:1883",
//			"topic":    queue,
//			"username": "",
//			"password": "",
//			"clientId": nuid.Next(),
//			"qos":      "0",
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
//func sendMessage(queue string, qos byte, data []byte) error {
//
//	opts := mqtt.NewClientOptions()
//	opts.AddBroker("tcp://localhost:1883")
//	opts.SetClientID(nuid.Next())
//	opts.SetConnectTimeout(defaultConnectTimeout)
//	client := mqtt.NewClient(opts)
//	if token := client.Connect(); token.Wait() && token.Error() != nil {
//		return fmt.Errorf("error connecting to mqtt broker, %w", token.Error())
//	}
//	token := client.Publish(queue, qos, false, data)
//	token.Wait()
//	return token.Error()
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
//		{
//			name: "request with target error",
//			target: &null.Client{
//				Delay:         0,
//				DoError:       fmt.Errorf("some-error"),
//				ResponseError: nil,
//			},
//			req:     types.NewRequest().SetData([]byte("some-data")),
//			queue:   "some-queue",
//			wantErr: false,
//		},
//		{
//			name:    "request with nil target",
//			target:  nil,
//			req:     types.NewRequest().SetData([]byte("some-data")),
//			queue:   "some-queue",
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//			defer cancel()
//			c, err := setupClient(ctx, tt.queue, tt.target)
//			if tt.wantErr {
//				require.Error(t, err)
//				return
//			}
//			require.NoError(t, err)
//			defer func() {
//				_ = c.Stop()
//			}()
//
//			err = sendMessage(tt.queue, 0, tt.req.Data)
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
//				Name: "mqtt",
//				Kind: "",
//				Properties: map[string]string{
//					"host":     "localhost:1883",
//					"topic":    "some-queue",
//					"username": "",
//					"password": "",
//					"clientId": nuid.Next(),
//					"qos":      "0",
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "init - bad url",
//			cfg: config.Spec{
//				Name: "mqtt-target",
//				Kind: "",
//				Properties: map[string]string{
//					"host":     "localhost:2000",
//					"topic":    "some-queue",
//					"username": "",
//					"password": "",
//					"clientId": nuid.Next(),
//					"qos":      "0",
//				},
//			},
//			wantErr: true,
//		},
//		{
//			name: "bad init - no  url",
//			cfg: config.Spec{
//				Name: "mqtt",
//				Kind: "",
//				Properties: map[string]string{
//					"topic":    "some-queue",
//					"username": "",
//					"password": "",
//					"clientId": nuid.Next(),
//					"qos":      "0",
//				},
//			},
//			wantErr: true,
//		},
//		{
//			name: "bad init - no topic",
//			cfg: config.Spec{
//				Name: "mqtt",
//				Kind: "",
//				Properties: map[string]string{
//					"host":     "localhost:1883",
//					"username": "",
//					"password": "",
//					"clientId": nuid.Next(),
//					"qos":      "0",
//				},
//			},
//			wantErr: true,
//		},
//		{
//			name: "init - bad qos",
//			cfg: config.Spec{
//				Name: "mqtt-target",
//				Kind: "",
//				Properties: map[string]string{
//					"host":     "localhost:1883",
//					"topic":    "some-queue",
//					"username": "",
//					"password": "",
//					"clientId": nuid.Next(),
//					"qos":      "-1",
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
//			require.EqualValues(t, tt.cfg.Name, c.Name())
//		})
//	}
//}
