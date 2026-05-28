package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Track total jobs processed (success vs failed)
	JobsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jobs_processed_total",
			Help: "Total processed jobs by status",
		},
		[]string{"status", "priority", "worker_id"},
	)

	// Track processing latency
	JobProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "job_processing_duration_seconds",
			Help:    "Time taken to process a job",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"priority"},
	)

	// Track enqueue rate
	JobsEnqueued = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jobs_enqueued_total",
			Help: "Total jobs enqueued into the system",
		},
		[]string{"priority"},
	)

	// Track jobs moving to Dead Letter Queue
	DeadLetterJobs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dead_letter_jobs_total",
			Help: "Total jobs moved to the DLQ",
		},
		[]string{"priority"},
	)
)