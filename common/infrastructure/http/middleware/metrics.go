package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics содержит RED-метрики + логические ошибки
type PrometheusMetrics struct {
	RequestCount   *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
	ServiceErrors  *prometheus.CounterVec
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
		ServiceErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "service_errors_total",
				Help: "Total number of service-level errors",
			},
			[]string{"type"},
		),
	}
}

// RegisterPrometheusMetrics регистрирует метрики в Prometheus
func (pm *PrometheusMetrics) Register() {
	prometheus.MustRegister(pm.RequestCount, pm.RequestLatency, pm.ServiceErrors)
}

// PrometheusMetricsMiddleware возвращает middleware для сбора метрик по RED + ошибок
func PrometheusMetricsMiddleware(metrics *PrometheusMetrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			// Собираем метрики
			latency := time.Since(start).Seconds()
			status := c.Response().Status
			statusStr := strconv.Itoa(status)

			metrics.RequestCount.WithLabelValues(c.Request().Method, c.Path(), statusStr).Inc()
			metrics.RequestLatency.WithLabelValues(c.Request().Method, c.Path()).Observe(latency)

			// Если произошла ошибка - классифицируем её
			if err != nil {
				errorType := classifyError(err)
				metrics.ServiceErrors.WithLabelValues(errorType).Inc()
			}

			// Добавляем обработку статусов 400-499 и 500+
			if status >= 400 {
				var errorType string
				if status >= 400 && status < 500 {
					errorType = "http_client"
				} else if status >= 500 {
					errorType = "http_server"
				}
				metrics.ServiceErrors.WithLabelValues(errorType).Inc()
			}

			return err
		}
	}
}

// classifyError - классифицирует ошибки для метрик
func classifyError(err error) string {
	// if errors.Is(err, domain.ErrDBConnection) {
	// 	return "database"
	// }
	// if errors.Is(err, domain.ErrTimeout) {
	// 	return "timeout"
	// }
	// if errors.Is(err, domain.ErrBusinessLogic) {
	// 	return "business_logic"
	// }
	return "unknown"
}
