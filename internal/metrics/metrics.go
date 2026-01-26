package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	HTTPRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	GoalsFetchCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "goals_fetch_total",
		Help: "Total number of goals fetch operations",
	}, []string{"status"})

	GoalsStoreCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "goals_store_total",
		Help: "Total number of goals store operations",
	}, []string{"status"})

	MirrorsPopulateCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mirrors_populate_total",
		Help: "Total number of mirrors populate operations",
	}, []string{"status"})

	RemoveOldGoalsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "remove_old_goals_total",
		Help: "Total number of remove old goals operations",
	}, []string{"status"})
)
