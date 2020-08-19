package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Report struct {
	Key            string  `json:"-"`
	Binding        string  `json:"binding"`
	SourceName     string  `json:"source_name"`
	SourceKind     string  `json:"source_kind"`
	TargetName     string  `json:"target_name"`
	TargetKind     string  `json:"target_kind"`
	RequestCount   float64 `json:"request_count"`
	RequestVolume  float64 `json:"request_volume"`
	ResponseCount  float64 `json:"response_count"`
	ResponseVolume float64 `json:"response_volume"`
	ErrorsCount    float64 `json:"errors_count"`
}

func (m *Report) labels() prometheus.Labels {
	return prometheus.Labels{
		"binding":     m.Binding,
		"source_name": m.SourceName,
		"source_kind": m.SourceKind,
		"target_name": m.TargetName,
		"target_kind": m.TargetKind,
	}
}

func (m *Report) Clone() *Report {
	return &Report{
		Key:            m.Key,
		Binding:        m.Binding,
		SourceName:     m.SourceName,
		SourceKind:     m.SourceKind,
		TargetName:     m.TargetName,
		TargetKind:     m.TargetKind,
		RequestCount:   m.RequestCount,
		RequestVolume:  m.RequestVolume,
		ResponseCount:  m.ResponseCount,
		ResponseVolume: m.ResponseVolume,
		ErrorsCount:    m.ErrorsCount,
	}
}
