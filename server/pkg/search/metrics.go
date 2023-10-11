package search

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/models"
	"go.uber.org/zap"
	"time"
)

const SlowQueryThresholdSeconds float64 = 1

var stdBuckets = []float64{0.1, 0.25, 0.5, 1, 2, 5, 10, 20}

type metrics struct {
	queryDurationSeconds      prometheus.Histogram
	predictionDurationSeconds prometheus.Histogram
	listTermsDurationSeconds  prometheus.Histogram
}

func newMetrics() *metrics {
	m := &metrics{
		queryDurationSeconds: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "search",
				Subsystem: "searcher",
				Name:      "query_duration_seconds",
				Help:      "Num seconds taken to execute search",
				Buckets:   stdBuckets,
			},
		),
		predictionDurationSeconds: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "search",
				Subsystem: "searcher",
				Name:      "prediction_duration_seconds",
				Help:      "Num seconds taken to execute query prediction",
				Buckets:   stdBuckets,
			},
		),
		listTermsDurationSeconds: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "search",
				Subsystem: "searcher",
				Name:      "list_terms_duration_seconds",
				Help:      "Num seconds taken to execute term list",
				Buckets:   stdBuckets,
			},
		),
	}
	prometheus.DefaultRegisterer.MustRegister(
		m.queryDurationSeconds,
		m.predictionDurationSeconds,
		m.listTermsDurationSeconds,
	)
	return m
}

type InstrumentedSearcher struct {
	s      Searcher
	m      *metrics
	logger *zap.Logger
}

func (i *InstrumentedSearcher) Search(ctx context.Context, f filter.Filter, page int32) (*api.SearchResultList, error) {
	startTime := time.Now().UnixMilli()
	defer func() {
		taken := float64(time.Now().UnixMilli()-startTime) / 1000
		i.m.queryDurationSeconds.Observe(taken)
		if taken > SlowQueryThresholdSeconds {
			i.logger.Warn("Slow search query detected", zap.String("filter", filter.MustPrint(f)), zap.Int32("page", page))
		}
	}()
	return i.s.Search(ctx, f, page)
}

func (i *InstrumentedSearcher) PredictSearchTerms(ctx context.Context, prefix string, exact bool, numPredictions int32, f filter.Filter) (*api.SearchTermPredictions, error) {
	startTime := time.Now().UnixMilli()
	defer func() {
		taken := float64(time.Now().UnixMilli()-startTime) / 1000
		i.m.predictionDurationSeconds.Observe(taken)
		if taken > SlowQueryThresholdSeconds {
			i.logger.Warn(
				"Slow predict query detected",
				zap.String("prefix", prefix),
				zap.Bool("exact", exact),
				zap.Int32("num_predictions", numPredictions),
				zap.String("filter", filter.MustPrint(f)),
			)
		}
	}()
	return i.s.PredictSearchTerms(ctx, prefix, exact, numPredictions, f)
}

func (i *InstrumentedSearcher) ListTerms(fieldName string, prefix string) (models.FieldValues, error) {
	startTime := time.Now().UnixMilli()
	defer func() {
		taken := float64(time.Now().UnixMilli()-startTime) / 1000
		i.m.listTermsDurationSeconds.Observe(taken)
		if taken > SlowQueryThresholdSeconds {
			i.logger.Warn(
				"Slow term list query detected",
				zap.String("field_name", fieldName),
				zap.String("prefix", prefix),
			)
		}
	}()
	return i.s.ListTerms(fieldName, prefix)
}

func InstrumentSearcher(s Searcher, logger *zap.Logger) Searcher {
	return &InstrumentedSearcher{s: s, m: newMetrics(), logger: logger}
}
