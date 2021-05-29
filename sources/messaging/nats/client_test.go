package nats

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/middleware"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type mockMiddleware struct {
	client      *kubemq.Client
	channelName string
}

func (m *mockMiddleware) Init() {

	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithCheckConnection(true),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))

	if err != nil {
		panic(err)
	}
	m.client = client
	m.channelName = "messaging.nats"
}

func (m *mockMiddleware) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	fmt.Println(request)
	r := types.NewResponse()
	r.SetData([]byte("ok"))
	r.SetMetadata(`"result":"ok"`)
	event := m.client.NewEvent()
	event.Channel = m.channelName
	event.Body = request.Data
	err := event.Send(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type testStructure struct {
	url                string
	subject            string
	username           string
	password           string
	token              string
	sslcertificatefile string
	sslcertificatekey  string
}

func getTestStructure() (*testStructure, error) {
	t := &testStructure{}
	dat, err := ioutil.ReadFile("./../../../credentials/nats/url.txt")
	if err != nil {
		return nil, err
	}
	t.url = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/subject.txt")
	if err != nil {
		return nil, err
	}
	t.subject = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/username.txt")
	if err != nil {
		return nil, err
	}
	t.username = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/password.txt")
	if err != nil {
		return nil, err
	}
	t.password = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/token.txt")
	if err != nil {
		return nil, err
	}
	t.token = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/certFile.pem")
	if err != nil {
		return nil, err
	}
	t.sslcertificatefile = string(dat)
	dat, err = ioutil.ReadFile("./../../../credentials/nats/certKey.pem")
	if err != nil {
		return nil, err
	}
	t.sslcertificatekey = string(dat)

	return t, nil
}

func TestClient_Init(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init -without tls",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":             dat.url,
					"subject":         dat.subject,
					"dynamic_mapping": "false",
					"username":        dat.username,
					"password":        dat.password,
					"token":           dat.password,
					"tls":             "false",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid init - no  url",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"subject":         dat.subject,
					"dynamic_mapping": "false",
					"username":        dat.username,
					"password":        dat.password,
					"token":           dat.password,
					"tls":             "true",
					"cert_file":       dat.sslcertificatefile,
					"cert_key":        dat.sslcertificatekey,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - missing cert key",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":             dat.url,
					"subject":         dat.subject,
					"dynamic_mapping": "false",
					"username":        dat.username,
					"password":        dat.password,
					"token":           dat.password,
					"tls":             "true",
					"cert_file":       dat.sslcertificatefile,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid init - missing cert file",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":             dat.url,
					"subject":         dat.subject,
					"dynamic_mapping": "false",
					"username":        dat.username,
					"password":        dat.password,
					"token":           dat.password,
					"tls":             "true",
					"cert_key":        dat.sslcertificatekey,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestClient_Do(t *testing.T) {
	dat, err := getTestStructure()
	require.NoError(t, err)
	middle := &mockMiddleware{}
	middle.Init()
	tests := []struct {
		name       string
		cfg        config.Spec
		wantErr    bool
		middleware middleware.Middleware
	}{
		{
			name: "start -with tls",
			cfg: config.Spec{
				Name: "messaging-nats",
				Kind: "messaging.nats",
				Properties: map[string]string{
					"url":             dat.url,
					"subject":         dat.subject,
					"dynamic_mapping": "false",
					"username":        dat.username,
					"password":        dat.password,
					"token":           dat.password,
					"tls":             "true",
					"cert_file":       dat.sslcertificatefile,
					"cert_key":        dat.sslcertificatekey,
				},
			},
			middleware: middle,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg, nil)
			require.NoError(t, err)
			err = c.Start(ctx, tt.middleware)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.Nil(t, err)
			time.Sleep(time.Duration(45) * time.Second)
			err = c.Stop()
			require.Nil(t, err)
			time.Sleep(time.Duration(10) * time.Second)
		})
	}
}
