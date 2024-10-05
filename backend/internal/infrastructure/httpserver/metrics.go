package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

// Создаем метрику для отслеживания количества горутин по эндпоинтам
var counterGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "http_active_goroutines",
		Help: "Number of active goroutines handling requests by endpoint",
	},
	[]string{"method", "endpoint"},
)

func init() {
	// Регистрируем метрику в Prometheus
	prometheus.MustRegister(counterGauge)
}

// Middleware для подсчета горутин
func counterMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Получаем путь и метод запроса
		endpoint := c.Path()         // Путь запроса (например, "/users")
		method := c.Request().Method // HTTP-метод (например, "GET")

		// Инкрементируем счетчик при начале обработки запроса
		counterGauge.WithLabelValues(method, endpoint).Inc()

		// Вызываем основной обработчик
		err := next(c)

		// Декрементируем счетчик после завершения запроса
		counterGauge.WithLabelValues(method, endpoint).Dec()

		// Возвращаем результат
		return err
	}
}
