package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type promCounterMetric struct {
	metric *prometheus.CounterVec
}

func newPromCounterMetric(subsystem, name, help string, labels ...string) *promCounterMetric {
	opts := prometheus.CounterOpts{
		Namespace:   "kubemq_sources",
		Subsystem:   subsystem,
		Name:        name,
		Help:        help,
		ConstLabels: nil,
	}

	c := &promCounterMetric{}
	c.metric = prometheus.NewCounterVec(opts, labels)
	return c
}

func (c *promCounterMetric) add(value float64, labels prometheus.Labels) {
	if value > 0 {
		c.metric.With(labels).Add(value)
	}

}
