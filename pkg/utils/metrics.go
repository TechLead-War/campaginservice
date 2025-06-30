package utils

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/cpu"
)

var (

	// Define your metrics here, e.g.:

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{
			"path",
			"method",
			"status_code",
		},
	)

	CacheActionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_actions_total",
			Help: "Total number of cache actions performed.",
		},
		[]string{"type"},
	)

	DBOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "Duration of database operations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// New metrics for performance comparison
	CPUGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "Current CPU usage percentage.",
		},
	)

	MemoryGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes.",
		},
	)

	GoroutineGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_total",
			Help: "Current number of goroutines.",
		},
	)

	// API specific latency metrics
	DeliveryAPILatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "delivery_api_latency_seconds",
			Help:    "Latency of delivery API endpoint.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
		},
	)
)

func InitMetrics() {
	// Start background goroutine to collect system metrics
	go collectSystemMetrics()
}

func RecordCacheHit() {
	CacheActionsTotal.With(prometheus.Labels{"type": "hit"}).Inc()
}

func RecordCacheMiss() {
	CacheActionsTotal.With(prometheus.Labels{"type": "miss"}).Inc()
}

// collectSystemMetrics collects CPU, memory, and goroutine metrics
func collectSystemMetrics() {
	ticker := time.NewTicker(15 * time.Second) // Collect every 15 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Record goroutine count
		GoroutineGauge.Set(float64(runtime.NumGoroutine()))

		// Record memory usage
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		MemoryGauge.Set(float64(m.Alloc)) // Allocated memory

		// Record CPU usage
		cpuPercentages, err := cpu.Percent(0, false)
		if err == nil && len(cpuPercentages) > 0 {
			CPUGauge.Set(cpuPercentages[0])
		}
	}
}
