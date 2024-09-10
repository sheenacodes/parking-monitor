package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	EventProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_processing_latency_seconds",
			Help:    "Latency of event processing in seconds.",
			Buckets: prometheus.DefBuckets, // Default buckets, can be customized
		},
		[]string{"event_type"},
	)

	EventProcessingFails = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_processing_fails_total",
			Help: "Total number of errors while processing events.",
		},
		[]string{"event_type", "error_stage"},
	)

	EventProcessingSuccesses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "event_processing_successes_total",
			Help: "Total number of errors while processing events.",
		},
		[]string{"event_type"},
	)
)

func init() {
	log.Println("register prometheus")
	prometheus.MustRegister(EventProcessingLatency)
	prometheus.MustRegister(EventProcessingFails)
	prometheus.MustRegister(EventProcessingSuccesses)
}
