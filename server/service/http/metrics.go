package http

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

type MetricsService struct {
}

func (c *MetricsService) RegisterHTTP(ctx context.Context, router *mux.Router) {
	router.Path("/metrics").Handler(promhttp.Handler())
}
