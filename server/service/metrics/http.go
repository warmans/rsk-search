package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewHTTPMetrics() *HTTPMetrics {
	metrics := &HTTPMetrics{
		OutboundMediaBytesTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "service",
				Subsystem:   "http",
				Name:        "outbound_bytes",
				Help:        "number of bytes written to response buffer",
				ConstLabels: nil,
			},
		),
		OutboundMediaQuotaRemaining: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "service",
				Subsystem:   "http",
				Name:        "outbound_bytes_quota_remaining",
				Help:        "remaining bytes in media quota",
				ConstLabels: nil,
			},
		),
	}
	prometheus.DefaultRegisterer.MustRegister(metrics.OutboundMediaBytesTotal)
	prometheus.DefaultRegisterer.MustRegister(metrics.OutboundMediaQuotaRemaining)
	return metrics
}

type HTTPMetrics struct {
	OutboundMediaBytesTotal     prometheus.Gauge
	OutboundMediaQuotaRemaining prometheus.Gauge
}
