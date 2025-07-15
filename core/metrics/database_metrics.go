package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// DatabaseMetrics holds all database-related metrics
type DatabaseMetrics struct {
	QueryDuration     *prometheus.HistogramVec
	QueryTotal        *prometheus.CounterVec
	QueryErrors       *prometheus.CounterVec
	ConnectionsActive prometheus.Gauge
	ViewRefreshes     *prometheus.CounterVec
	ViewSize          *prometheus.GaugeVec
	ViewRowCount      *prometheus.GaugeVec
}

// NewDatabaseMetrics creates and registers database metrics
func NewDatabaseMetrics() *DatabaseMetrics {
	return &DatabaseMetrics{
		QueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "analabit_database_query_duration_seconds",
				Help: "Duration of database queries in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"query_type", "table", "status"},
		),
		QueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "analabit_database_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"query_type", "table", "status"},
		),
		QueryErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "analabit_database_query_errors_total",
				Help: "Total number of database query errors",
			},
			[]string{"query_type", "table", "error_type"},
		),
		ConnectionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "analabit_database_connections_active",
				Help: "Number of active database connections",
			},
		),
		ViewRefreshes: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "analabit_materialized_view_refreshes_total",
				Help: "Total number of materialized view refreshes",
			},
			[]string{"view_name", "refresh_type", "status"},
		),
		ViewSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "analabit_materialized_view_size_bytes",
				Help: "Size of materialized views in bytes",
			},
			[]string{"view_name"},
		),
		ViewRowCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "analabit_materialized_view_rows",
				Help: "Number of rows in materialized views",
			},
			[]string{"view_name"},
		),
	}
}

// RecordQuery records metrics for a database query
func (m *DatabaseMetrics) RecordQuery(queryType, table string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
		m.QueryErrors.WithLabelValues(queryType, table, "database_error").Inc()
	}

	m.QueryDuration.WithLabelValues(queryType, table, status).Observe(duration.Seconds())
	m.QueryTotal.WithLabelValues(queryType, table, status).Inc()
}

// RecordViewRefresh records metrics for materialized view refresh
func (m *DatabaseMetrics) RecordViewRefresh(viewName, refreshType string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	m.ViewRefreshes.WithLabelValues(viewName, refreshType, status).Inc()
	m.QueryDuration.WithLabelValues("view_refresh", viewName, status).Observe(duration.Seconds())
}

// UpdateViewStats updates materialized view statistics
func (m *DatabaseMetrics) UpdateViewStats(viewName string, size int64, rowCount int64) {
	m.ViewSize.WithLabelValues(viewName).Set(float64(size))
	m.ViewRowCount.WithLabelValues(viewName).Set(float64(rowCount))
}

// SetActiveConnections sets the number of active database connections
func (m *DatabaseMetrics) SetActiveConnections(count int) {
	m.ConnectionsActive.Set(float64(count))
}

// Global metrics instance
var DBMetrics *DatabaseMetrics

// InitMetrics initializes the global metrics instance
func InitMetrics() {
	DBMetrics = NewDatabaseMetrics()
}