package middleware

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-sources/config"
	"github.com/kubemq-hub/kubemq-sources/pkg/metrics"
)

type MetricsMiddleware struct {
	exporter     *metrics.Exporter
	metricReport *metrics.Report
}

func NewMetricsMiddleware(cfg config.BindingConfig, exporter *metrics.Exporter) (*MetricsMiddleware, error) {
	if exporter == nil {
		return nil, fmt.Errorf("no valid exporter found")
	}
	m := &MetricsMiddleware{
		exporter: exporter,
		metricReport: &metrics.Report{
			Key:            fmt.Sprintf("%s-%s-%s", cfg.Name, cfg.Source.Kind, cfg.Target.Kind),
			Binding:        cfg.Name,
			SourceKind:     cfg.Source.Kind,
			TargetKind:     cfg.Target.Kind,
			RequestCount:   0,
			RequestVolume:  0,
			ResponseCount:  0,
			ResponseVolume: 0,
			ErrorsCount:    0,
		},
	}
	return m, nil
}

func (m *MetricsMiddleware) clearReport() {
	m.metricReport.ErrorsCount = 0
	m.metricReport.ResponseVolume = 0
	m.metricReport.ResponseCount = 0
	m.metricReport.RequestVolume = 0
	m.metricReport.RequestCount = 0
}
