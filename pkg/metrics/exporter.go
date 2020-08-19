package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var labels = []string{"binding", "source_name", "source_kind", "target_name", "target_kind"}

type Exporter struct {
	Store                    *Store
	requestsCollector        *promCounterMetric
	responsesCollector       *promCounterMetric
	requestsVolumeCollector  *promCounterMetric
	responsesVolumeCollector *promCounterMetric
	errorsCollector          *promCounterMetric
}

func (e *Exporter) PrometheusHandler() http.Handler {
	return promhttp.Handler()
}

func NewExporter() (*Exporter, error) {
	e := &Exporter{
		Store:                    NewStore(),
		requestsCollector:        nil,
		responsesCollector:       nil,
		requestsVolumeCollector:  nil,
		responsesVolumeCollector: nil,
		errorsCollector:          nil,
	}
	if err := e.initPromMetrics(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Exporter) initPromMetrics() error {
	e.requestsCollector = newPromCounterMetric(
		"requests",
		"count",
		"counts requests per binding,source and target types",
		labels...,
	)
	e.responsesCollector = newPromCounterMetric(
		"responses",
		"count",
		"counts responses per binding,source and target types",
		labels...,
	)
	e.requestsVolumeCollector = newPromCounterMetric(
		"requests",
		"volume",
		"sum requests volume per binding,source and target types",
		labels...,
	)
	e.responsesVolumeCollector = newPromCounterMetric(
		"responses",
		"volume",
		"sum responses volume per binding,source and target types",
		labels...,
	)
	e.errorsCollector = newPromCounterMetric(
		"errors",
		"count",
		"counts error requests per binding,source and target types",
		labels...,
	)

	err := prometheus.Register(e.requestsCollector.metric)
	if err != nil {
		return err
	}
	err = prometheus.Register(e.responsesCollector.metric)
	if err != nil {
		return err
	}
	err = prometheus.Register(e.requestsVolumeCollector.metric)
	if err != nil {
		return err
	}
	err = prometheus.Register(e.responsesVolumeCollector.metric)
	if err != nil {
		return err
	}
	err = prometheus.Register(e.errorsCollector.metric)
	if err != nil {
		return err
	}

	return nil
}

func (e *Exporter) Report(m *Report) {
	lbs := m.labels()
	e.requestsCollector.add(m.RequestCount, lbs)
	e.requestsVolumeCollector.add(m.RequestVolume, lbs)
	e.responsesCollector.add(m.ResponseCount, lbs)
	e.responsesVolumeCollector.add(m.ResponseVolume, lbs)
	e.errorsCollector.add(m.ErrorsCount, lbs)
	e.Store.Add(m)
}
