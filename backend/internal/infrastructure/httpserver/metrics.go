package httpserver

import (
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// httpMetrics — структура для хранения метрик
type httpMetrics struct {
	counterGauge    *prometheus.GaugeVec
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.SummaryVec
	responseSize    *prometheus.SummaryVec
}

// NewMetricsMiddleware создает мидлварь для сбора метрик с переданными параметрами
func NewMetricsMiddleware(registry prometheus.Registerer, cfg config.MetricsConfig) echo.MiddlewareFunc {

	bucketsHttpRequestDuration := cfg.BucketsHttpRequestDuration
	if len(bucketsHttpRequestDuration) == 0 {
		bucketsHttpRequestDuration = prometheus.DefBuckets
	}

	// Создаем метрики
	httpMetrics := &httpMetrics{
		counterGauge: promauto.With(registry).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_active_goroutines",
				Help: "Number of active goroutines handling requests by endpoint",
			},
			[]string{"method", "endpoint"},
		),
		requestsTotal: promauto.With(registry).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Tracks the number of HTTP requests.",
			},
			[]string{"method", "code", "endpoint"},
		),
		requestDuration: promauto.With(registry).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Tracks the latencies for HTTP requests.",
				Buckets: bucketsHttpRequestDuration,
			},
			[]string{"method", "code", "endpoint"},
		),
		requestSize: promauto.With(registry).NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_request_size_bytes",
				Help: "Tracks the size of HTTP requests.",
			},
			[]string{"method", "code", "endpoint"},
		),
		responseSize: promauto.With(registry).NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_response_size_bytes",
				Help: "Tracks the size of HTTP responses.",
			},
			[]string{"method", "code", "endpoint"},
		),
	}

	// Возвращаем функцию мидлвари, которая использует метрики из конфигурации
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			endpoint := c.Path()
			method := c.Request().Method

			// Подсчитываем активные горутины
			httpMetrics.counterGauge.WithLabelValues(method, endpoint).Inc()
			defer httpMetrics.counterGauge.WithLabelValues(method, endpoint).Dec()

			// Запускаем основной обработчик
			err := next(c)

			// Завершаем отслеживание времени
			duration := time.Since(start).Seconds()
			code := fmt.Sprintf("%d", c.Response().Status)

			// Обновляем метрики
			httpMetrics.requestsTotal.WithLabelValues(method, code, endpoint).Inc()
			httpMetrics.requestDuration.WithLabelValues(method, code, endpoint).Observe(duration)
			httpMetrics.requestSize.WithLabelValues(method, code, endpoint).Observe(float64(c.Request().ContentLength))
			httpMetrics.responseSize.WithLabelValues(method, code, endpoint).Observe(float64(c.Response().Size))

			return err
		}
	}
}
