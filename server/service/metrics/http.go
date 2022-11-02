package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewHTTPMetrics() *HTTPMetrics {
	metrics := &HTTPMetrics{
		OutboundMediaBytesTotal: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:   "service",
				Subsystem:   "http",
				Name:        "outbound_bytes_total",
				Help:        "sum of bytes written to response buffer",
				ConstLabels: nil,
				Objectives:  nil,
				MaxAge:      0,
				AgeBuckets:  0,
				BufCap:      0,
			},
			[]string{"media_type"},
		),
	}
	prometheus.DefaultRegisterer.Register(metrics.OutboundMediaBytesTotal)
	return metrics
}

type HTTPMetrics struct {
	OutboundMediaBytesTotal *prometheus.SummaryVec
}
