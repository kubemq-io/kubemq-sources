package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/types"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockTarget struct {
	setTargetErr    error
	setExecutionErr error
}

func (m *mockTarget) Init(ctx context.Context, cfg config.Spec) error {
	return nil
}

func (m *mockTarget) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	if m.setTargetErr != nil {
		return nil, m.setTargetErr
	}
	if m.setExecutionErr != nil {
		return types.NewResponse().SetError(m.setExecutionErr), nil
	}
	return types.NewResponse().SetData(request.Data), nil
}

func TestHandler_ServeHTTP_POST(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	h := New()
	_ = h.Init(ctx, config.Spec{
		Name: "",
		Kind: "",
		Properties: map[string]string{
			"methods": "post",
			"path":    "/",
		},
	}, nil)
	err := h.Start(ctx, &mockTarget{})
	require.NoError(t, err)
	defer func() {
		_ = h.Stop()
	}()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ServeHTTP)
	req, err := http.NewRequestWithContext(ctx, "POST", "/do", bytes.NewBufferString("some-data"))
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)
	require.Equal(t, 200, rr.Code)
	body, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.EqualValues(t, "some-data", string(body))
}

func TestHandler_ServeHTTP_GET(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	h := New()
	_ = h.Init(ctx, config.Spec{
		Name: "",
		Kind: "",
		Properties: map[string]string{
			"methods": "get",
			"path":    "/",
		},
	}, nil)
	err := h.Start(context.Background(), &mockTarget{})
	require.NoError(t, err)
	defer func() {
		_ = h.Stop()
	}()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ServeHTTP)
	req, err := http.NewRequestWithContext(ctx, "GET", "/do", nil)
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)
	require.Equal(t, 200, rr.Code)
}
func TestHandler_ServeHTTP_TargetError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := New()

	_ = h.Init(ctx, config.Spec{
		Name: "",
		Kind: "",
		Properties: map[string]string{
			"methods": "get",
			"path":    "/",
		},
	}, nil)
	err := h.Start(context.Background(), &mockTarget{
		setTargetErr:    fmt.Errorf("error"),
		setExecutionErr: nil,
	})
	require.NoError(t, err)
	defer func() {
		_ = h.Stop()
	}()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ServeHTTP)
	req, err := http.NewRequestWithContext(ctx, "POST", "/do", bytes.NewBufferString("some-data"))
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)
	require.Equal(t, 500, rr.Code)
	body, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.EqualValues(t, "error\n", string(body))
}

func TestHandler_ServeHTTP_ExecutionError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	h := New()

	_ = h.Init(ctx, config.Spec{
		Name: "",
		Kind: "",
		Properties: map[string]string{
			"methods": "get",
			"path":    "/",
		},
	}, nil)
	err := h.Start(context.Background(), &mockTarget{
		setTargetErr:    nil,
		setExecutionErr: fmt.Errorf("error"),
	})
	require.NoError(t, err)
	defer func() {
		_ = h.Stop()
	}()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ServeHTTP)
	req, err := http.NewRequestWithContext(ctx, "POST", "/do", bytes.NewBufferString("some-data"))
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)
	require.Equal(t, 500, rr.Code)
	body, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	require.EqualValues(t, "error\n", string(body))
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "http-source",
				Kind: "",
				Properties: map[string]string{
					"methods": "post,get",
					"path":    "/",
				},
			},
			wantErr: false,
		},
		{
			name: "init - no methods",
			cfg: config.Spec{
				Name: "http-source",
				Kind: "",
				Properties: map[string]string{
					"path": "/",
				},
			},
			wantErr: true,
		},
		{
			name: "init - no path",
			cfg: config.Spec{
				Name: "http-source",
				Kind: "",
				Properties: map[string]string{
					"methods": "post,get",
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

			if err := c.Init(ctx, tt.cfg, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
