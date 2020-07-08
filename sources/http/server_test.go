package http

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kubemq-hub/kubemq-source-connectors/config"
	"github.com/kubemq-hub/kubemq-source-connectors/targets/null"
	"github.com/kubemq-hub/kubemq-source-connectors/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	targetMiddleware "github.com/kubemq-hub/kubemq-source-connectors/middleware"
)

func sendRequest(ctx context.Context, req *types.Request, url string) (*types.Response, error) {
	resp := &types.Response{}
	_, err := resty.New().R().SetBody(req).SetResult(resp).Post(url)
	return resp, err
}
func TestServer_process(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.Metadata
		target   targetMiddleware.Middleware
		req      *types.Request
		wantResp *types.Response
		url      string
		wantErr  bool
	}{
		{
			name: "request",
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "40000",
					"path": "/",
				},
			},
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse().SetData([]byte("some-data")),
			url:      "http://localhost:40000/",
			wantErr:  false,
		},
		{
			name: "request - target error",
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "40001",
					"path": "/",
				},
			},
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: fmt.Errorf("error"),
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse().SetMetadataKeyValue("error", "error"),
			url:      "http://localhost:40001/",
			wantErr:  false,
		},
		{
			name: "request - target error 2",
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "40002",
					"path": "/",
				},
			},
			target: &null.Client{
				Delay:         0,
				DoError:       fmt.Errorf("error"),
				ResponseError: nil,
			},
			req:      types.NewRequest().SetData([]byte("some-data")),
			wantResp: types.NewResponse().SetMetadataKeyValue("error", "error"),
			url:      "http://localhost:40002/",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			s := New()
			_ = s.Init(ctx, tt.cfg)
			err := s.Start(ctx, tt.target)
			defer func() {
				_ = s.Stop()
			}()
			require.NoError(t, err)
			gotResp, err := sendRequest(ctx, tt.req, tt.url)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
		})
	}
}

func TestServer_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "40000",
					"path": "",
				},
			},
			wantErr: false,
		},
		{
			name: "init",
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "-1",
					"path": "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}

func TestClient_Start(t *testing.T) {

	tests := []struct {
		name    string
		target  targetMiddleware.Middleware
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "start",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "41000",
					"path": "",
				},
			},
			wantErr: false,
		},
		{
			name:   "start - no target",
			target: nil,
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "",
					"port": "41000",
					"path": "",
				},
			},
			wantErr: true,
		},
		{
			name: "start - bad host",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			cfg: config.Metadata{
				Name: "http",
				Kind: "",
				Properties: map[string]string{
					"host": "12.7....0",
					"port": "41001",
					"path": "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			_ = c.Init(ctx, tt.cfg)

			if err := c.Start(ctx, tt.target); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
