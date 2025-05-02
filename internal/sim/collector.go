package sim

import (
	"fmt"
	"net/http"
	"log"
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PayloadMonitorCollector implements the prometheus.Collector interface
type PayloadMonitorCollector struct {
	monitor *PayloadMonitor // Pointer to the application's monitor state

	// Define descriptions for each metric we want to expose
	numPayloads *prometheus.Desc
	numTx       *prometheus.Desc
	numRx       *prometheus.Desc
	numErr      *prometheus.Desc
	bytesTx     *prometheus.Desc // Changed from Mbs to Bytes for convention
	bytesRx     *prometheus.Desc // Changed from Mbs to Bytes
	errBytes    *prometheus.Desc // Changed from Mbs to Bytes
	totalProcessingTimeSeconds *prometheus.Desc // Total time
	processedOperations *prometheus.Desc // Total operations processed
}

// NewPayloadMonitorCollector creates a new collector for the given monitor.
func NewPayloadMonitorCollector(monitor *PayloadMonitor) *PayloadMonitorCollector {
	return &PayloadMonitorCollector{
		monitor: monitor,
		numPayloads: prometheus.NewDesc(
			"app_payload_monitor_current_payloads", // Metric name (lowercase_underscore)
			"Current number of active payloads being monitored.", // Help text
			nil, // Label names (none for these simple metrics)
			nil, // Constant labels (none)
		),
		numTx: prometheus.NewDesc(
			"app_payload_monitor_transmissions_total", // Convention: _total for cumulative counters
			"Total number of transmissions.",
			nil,
			nil,
		),
		numRx: prometheus.NewDesc(
			"app_payload_monitor_receptions_total", // Convention: _total
			"Total number of receptions.",
			nil,
			nil,
		),
		numErr: prometheus.NewDesc(
			"app_payload_monitor_errors_total", // Convention: _total
			"Total number of errors encountered.",
			nil,
			nil,
		),
		bytesTx: prometheus.NewDesc(
			"app_payload_monitor_bytes_transmitted_total", // Convention: bytes and _total
			"Total number of bytes transmitted.",
			nil,
			nil,
		),
		bytesRx: prometheus.NewDesc(
			"app_payload_monitor_bytes_received_total", // Convention: bytes and _total
			"Total number of bytes received.",
			nil,
			nil,
		),
		errBytes: prometheus.NewDesc(
			"app_payload_monitor_error_bytes_total", // Convention: bytes and _total
			"Total number of bytes associated with errors.",
			nil,
			nil,
		),
		totalProcessingTimeSeconds: prometheus.NewDesc(
            "app_payload_monitor_processing_time_seconds_total", // Convention: seconds and _total
            "Total cumulative processing time in seconds.",
            nil,
            nil,
        ),
        processedOperations: prometheus.NewDesc(
            "app_payload_monitor_processed_operations_total", // Convention: _total
            "Total number of operations whose processing time is recorded.",
            nil,
            nil,
        ),
	}
}

// Describe sends the metric descriptions to the provided channel.
// This is called by the Prometheus client library during initialization.
func (collector *PayloadMonitorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.numPayloads
	ch <- collector.numTx
	ch <- collector.numRx
	ch <- collector.numErr
	ch <- collector.bytesTx
	ch <- collector.bytesRx
	ch <- collector.errBytes
	ch <- collector.totalProcessingTimeSeconds
    ch <- collector.processedOperations
}

// Collect reads the current state and sends metrics to the provided channel.
// This is called by the Prometheus client library whenever the /metrics endpoint is scraped.
func (collector *PayloadMonitorCollector) Collect(ch chan<- prometheus.Metric) {

	// Safely access the data from the monitor
	collector.monitor.mu.RLock() // Use RLock for reading
	numPayloads := collector.monitor.numPayloads
	numTx := collector.monitor.numTx
	numRx := collector.monitor.numRx
	numErr := collector.monitor.numErr

	// Convert MB back to Bytes for the metric value
	bytesTx := collector.monitor.numMbsTx * 1024 * 1024
	bytesRx := collector.monitor.numMbsRx * 1024 * 1024
	errBytes := collector.monitor.numErrMbs * 1024 * 1024
	totalProcessingTimeSeconds := collector.monitor.totalProcessingTime.Seconds() // Convert Duration to seconds
	processedOperations := collector.monitor.processedOperations
	collector.monitor.mu.RUnlock()

	// Send the collected data as Prometheus metrics
	ch <- prometheus.MustNewConstMetric(
		collector.numPayloads,
		prometheus.GaugeValue, // Gauge is for values that can go up or down
		float64(numPayloads),
	)
	ch <- prometheus.MustNewConstMetric(
		collector.numTx,
		prometheus.CounterValue, // Counter is for values that only increase
		float64(numTx),
	)
	ch <- prometheus.MustNewConstMetric(
		collector.numRx,
		prometheus.CounterValue,
		float64(numRx),
	)
	ch <- prometheus.MustNewConstMetric(
		collector.numErr,
		prometheus.CounterValue,
		float64(numErr),
	)
	ch <- prometheus.MustNewConstMetric(
		collector.bytesTx,
		prometheus.CounterValue,
		bytesTx, // Value is already float64
	)
	ch <- prometheus.MustNewConstMetric(
		collector.bytesRx,
		prometheus.CounterValue,
		bytesRx, // Value is already float64
	)
	ch <- prometheus.MustNewConstMetric(
		collector.errBytes,
		prometheus.CounterValue,
		errBytes, // Value is already float64
	)
	ch <- prometheus.MustNewConstMetric(
         collector.totalProcessingTimeSeconds,
         prometheus.CounterValue, // Total time is cumulative
         totalProcessingTimeSeconds,
     )
     ch <- prometheus.MustNewConstMetric(
         collector.processedOperations,
         prometheus.CounterValue, // Total operations is cumulative
         float64(processedOperations),
     )

    // Note: To get the average processing time in Grafana, you'd query
    // `rate(app_payload_monitor_processing_time_seconds_total[5m]) / rate(app_payload_monitor_processed_operations_total[5m])`
    // (adjusting the time window [5m] as needed)
}
// Initialize the Prometheus HTTP handler
func InitMonitoring(cfg *Config, infoChan <-chan PacketInfo) {
	// Create application's monitor instance
	monitor := NewPayloadMonitor()

	// Start payload monitor
	go procPayloadMon(cfg, monitor, infoChan)

	// Create the Prometheus collector for monitor
	collector := NewPayloadMonitorCollector(monitor)

	// Create a Prometheus registry and register collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Set up the HTTP server to expose metrics
	//    Use promhttp.HandlerFor with custom registry
	http.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			// Optional: enable metrics about the handler itself
			// EnableRequestCompression: true,
		},
	))

	// Start the HTTP server
	listenAddr := ":8080"
	if cfg.MetricsPort != 0 {
		listenAddr = fmt.Sprintf(":%d", cfg.MetricsPort)
	}
	fmt.Printf("Serving metrics on http://localhost%s/metrics\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}