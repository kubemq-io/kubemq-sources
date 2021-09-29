package middleware

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-sources/config"
	"github.com/kubemq-io/kubemq-sources/pkg/logger"
	"github.com/kubemq-io/kubemq-sources/pkg/metrics"
	"github.com/kubemq-io/kubemq-sources/types"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"
)

type mockTarget struct {
	request  *types.Request
	response *types.Response
	err      error
	delay    time.Duration
	executed int
}

func (m *mockTarget) Init(ctx context.Context, cfg config.Spec) error {
	return nil
}

func (m *mockTarget) Do(ctx context.Context, request *types.Request) (*types.Response, error) {
	time.Sleep(m.delay)
	m.executed++
	return m.response, m.err
}

func (m *mockTarget) Name() string {
	return ""
}

func TestClient_RateLimiter(t *testing.T) {
	tests := []struct {
		name             string
		mock             *mockTarget
		meta             types.Metadata
		timeToRun        time.Duration
		wantMaxExecution int
		wantErr          bool
	}{
		{
			name: "100 per seconds",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      nil,
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"rate_per_seconds": "100",
			},
			timeToRun:        time.Second,
			wantMaxExecution: 110,
			wantErr:          false,
		},
		{
			name: "unlimited",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      nil,
				delay:    0,
				executed: 0,
			},
			meta:             map[string]string{},
			timeToRun:        time.Second,
			wantMaxExecution: math.MaxInt32,
			wantErr:          false,
		},
		{
			name: "bad rpc",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      nil,
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"rate_per_seconds": "-100",
			},
			timeToRun:        time.Second,
			wantMaxExecution: 0,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeToRun)
			defer cancel()
			rl, err := NewRateLimitMiddleware(tt.meta)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			md := Chain(tt.mock, RateLimiter(rl))
			for {
				select {
				case <-ctx.Done():
					goto done
				default:

				}
				_, _ = md.Do(ctx, tt.mock.request)
			}
		done:
			require.LessOrEqual(t, tt.mock.executed, tt.wantMaxExecution)

		})
	}
}

func TestClient_Retry(t *testing.T) {
	log := logger.NewLogger("TestClient_Retry")
	tests := []struct {
		name        string
		mock        *mockTarget
		meta        types.Metadata
		wantRetries int
		wantErr     bool
	}{
		{
			name: "no retry options",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta:        map[string]string{},
			wantRetries: 1,
			wantErr:     false,
		},
		{
			name: "retry with success",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      nil,
				delay:    0,
				executed: 0,
			},
			meta:        map[string]string{},
			wantRetries: 1,
			wantErr:     false,
		},
		{
			name: "3 retries fixed delay",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":           "3",
				"retry_delay_milliseconds": "100",
				"retry_delay_type":         "fixed",
			},
			wantRetries: 3,
			wantErr:     false,
		},
		{
			name: "3 retries back-off delay",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":           "3",
				"retry_delay_milliseconds": "100",
				"retry_delay_type":         "back-off",
			},
			wantRetries: 3,
			wantErr:     false,
		},
		{
			name: "3 retries random delay",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":           "3",
				"retry_delay_milliseconds": "200",
				"retry_delay_type":         "random",
			},
			wantRetries: 3,
			wantErr:     false,
		},
		{
			name: "3 retries random with jitter delay",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":                "3",
				"retry_delay_milliseconds":      "200",
				"retry_max_jitter_milliseconds": "200",
				"retry_delay_type":              "random",
			},
			wantRetries: 3,
			wantErr:     false,
		},
		{
			name: "bad rate limiter - attempts",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":                "-3",
				"retry_delay_milliseconds":      "200",
				"retry_max_jitter_milliseconds": "200",
				"retry_delay_type":              "random",
			},
			wantRetries: 3,
			wantErr:     true,
		},
		{
			name: "bad rate limiter - retry_delay_milliseconds",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":                "3",
				"retry_delay_milliseconds":      "-200",
				"retry_max_jitter_milliseconds": "200",
				"retry_delay_type":              "random",
			},
			wantRetries: 3,
			wantErr:     true,
		},
		{
			name: "bad rate limiter - retry_max_jitter_milliseconds",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":                "3",
				"retry_delay_milliseconds":      "200",
				"retry_max_jitter_milliseconds": "-200",
				"retry_delay_type":              "random",
			},
			wantRetries: 3,
			wantErr:     true,
		},
		{
			name: "bad rate limiter - retry_delay_type",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"retry_attempts":                "3",
				"retry_delay_milliseconds":      "200",
				"retry_max_jitter_milliseconds": "200",
				"retry_delay_type":              "bad-type",
			},
			wantRetries: 3,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			r, err := NewRetryMiddleware(tt.meta, log)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			md := Chain(tt.mock, Retry(r))
			resp, err := md.Do(ctx, tt.mock.request)
			if tt.mock.err != nil {
				require.Error(t, err)
				require.EqualValues(t, tt.wantRetries, tt.mock.executed)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

		})
	}
}
func TestClient_Metric(t *testing.T) {
	exporter, err := metrics.NewExporter()
	require.NoError(t, err)

	tests := []struct {
		name       string
		mock       *mockTarget
		cfg        config.BindingConfig
		wantReport *metrics.Report
		wantErr    bool
	}{
		{
			name: "no error request",
			mock: &mockTarget{
				request:  types.NewRequest().SetData([]byte("data")),
				response: types.NewResponse().SetData([]byte("data")),
				err:      nil,
				delay:    0,
				executed: 0,
			},
			cfg: config.BindingConfig{
				Name: "b-1",
				Source: config.Spec{
					Name:       "sn",
					Kind:       "sk",
					Properties: nil,
				},
				Target: config.Spec{
					Name:       "tn",
					Kind:       "tk",
					Properties: nil,
				},
				Properties: nil,
			},
			wantReport: &metrics.Report{
				Key:            "b-1-sn-sk-tn-tk",
				Binding:        "b-1",
				SourceKind:     "sk",
				TargetKind:     "tk",
				RequestCount:   1,
				RequestVolume:  4,
				ResponseCount:  1,
				ResponseVolume: 4,
				ErrorsCount:    0,
			},
			wantErr: false,
		},
		{
			name: "error request",
			mock: &mockTarget{
				request:  types.NewRequest().SetData([]byte("data")),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			cfg: config.BindingConfig{
				Name: "b-2",
				Source: config.Spec{
					Name:       "sn",
					Kind:       "sk",
					Properties: nil,
				},
				Target: config.Spec{
					Name:       "tn",
					Kind:       "tk",
					Properties: nil,
				},
				Properties: nil,
			},
			wantReport: &metrics.Report{
				Key:            "b-2-sn-sk-tn-tk",
				Binding:        "b-2",
				SourceKind:     "sk",
				TargetKind:     "tk",
				RequestCount:   1,
				RequestVolume:  4,
				ResponseCount:  0,
				ResponseVolume: 0,
				ErrorsCount:    1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			m, err := NewMetricsMiddleware(tt.cfg, exporter)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			md := Chain(tt.mock, Metric(m))
			_, _ = md.Do(ctx, tt.mock.request)
			storedReport := exporter.Store.Get(tt.wantReport.Key)
			require.EqualValues(t, tt.wantReport, storedReport)
		})
	}
}
func TestClient_Log(t *testing.T) {

	tests := []struct {
		name    string
		mock    *mockTarget
		meta    types.Metadata
		wantErr bool
	}{
		{
			name: "no log level",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: nil,
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta:    map[string]string{},
			wantErr: false,
		},
		{
			name: "debug level",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"log_level": "debug",
			},
			wantErr: false,
		},
		{
			name: "info level",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"log_level": "info",
			},
			wantErr: false,
		},
		{
			name: "info level - 2",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      nil,
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"log_level": "info",
			},
			wantErr: false,
		},
		{
			name: "error level",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"log_level": "error",
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			mock: &mockTarget{
				request:  types.NewRequest(),
				response: types.NewResponse(),
				err:      fmt.Errorf("some-error"),
				delay:    0,
				executed: 0,
			},
			meta: map[string]string{
				"log_level": "bad-level",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			log, err := NewLogMiddleware("test", tt.meta)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			md := Chain(tt.mock, Log(log))
			_, _ = md.Do(ctx, tt.mock.request)
		})
	}
}
func TestClient_Chain(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mock := &mockTarget{
		request:  types.NewRequest(),
		response: nil,
		err:      fmt.Errorf("some-error"),
		delay:    0,
		executed: 0,
	}
	meta := map[string]string{
		"log_level":                     "debug",
		"rate_per_seconds":              "1",
		"retry_attempts":                "3",
		"retry_delay_milliseconds":      "100",
		"retry_max_jitter_milliseconds": "100",
		"retry_delay_type":              "fixed",
	}
	log, err := NewLogMiddleware("test", meta)
	require.NoError(t, err)
	rl, err := NewRateLimitMiddleware(meta)
	require.NoError(t, err)
	retry, err := NewRetryMiddleware(meta, logger.NewLogger("retry-logger"))
	require.NoError(t, err)
	md := Chain(mock, RateLimiter(rl), Retry(retry), Log(log))
	start := time.Now()
	_, _ = md.Do(ctx, mock.request)
	d := time.Since(start)
	require.GreaterOrEqual(t, d.Milliseconds(), 2*time.Second.Milliseconds())
}
