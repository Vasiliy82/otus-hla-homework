package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics содержит метрики
type PrometheusMetrics struct {
	RequestCount   *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
}

// NewPrometheusMetrics создает объект метрик
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		RequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Histogram of response latency (seconds)",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}

// RegisterPrometheusMetrics регистрирует метрики в Prometheus
func (pm *PrometheusMetrics) Register() {
	prometheus.MustRegister(pm.RequestCount, pm.RequestLatency)
}

// PrometheusMetricsMiddleware возвращает middleware для сбора метрик
func PrometheusMetricsMiddleware(metrics *PrometheusMetrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Выполняем обработку запроса
			err := next(c)

			// Собираем метрики
			latency := time.Since(start).Seconds()
			status := c.Response().Status
			metrics.RequestCount.WithLabelValues(c.Request().Method, c.Path(), string(rune(status))).Inc()
			metrics.RequestLatency.WithLabelValues(c.Request().Method, c.Path()).Observe(latency)

			return err
		}
	}
}
