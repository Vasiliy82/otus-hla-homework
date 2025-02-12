package rest

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/labstack/echo/v4"
)

const (
	prefix                 = "application/"
	defaultProtocolVersion = "default"
)

type ProxyRouter struct {
	proxies map[string]map[string]*httputil.ReverseProxy
}

// NewProxyRouter инициализирует все прокси на основе конфигурации
func NewProxyRouter(cfg *config.SocialNetworkConfig) (*ProxyRouter, error) {
	proxies := make(map[string]map[string]*httputil.ReverseProxy)

	for _, route := range cfg.RoutingConfig {
		proxies[route.Path] = make(map[string]*httputil.ReverseProxy)
		for _, service := range route.Services {
			if _, exists := proxies[route.Path][service.URL]; !exists {
				target, err := url.Parse(service.URL)
				if err != nil {
					return nil, err
				}

				newproxy := httputil.NewSingleHostReverseProxy(target)

				newproxy.ModifyResponse = func(resp *http.Response) error {
					requestID := resp.Header.Get("X-Request-ID")
					if requestID != "" {
						// Если заголовок уже есть, удаляем
						resp.Header.Del("X-Request-ID")
					}
					return nil
				}

				proxies[route.Path][service.URL] = newproxy

				// Здесь могут быть настроены всевозможные кастомные параметры транспорта: таймауты,
				// idle time, keepalive, и прочее.
				// proxies[route.Path][service.URL].Transport = custom_transport
			}
		}
	}

	return &ProxyRouter{proxies: proxies}, nil
}

// RouterHandler обрабатывает маршруты в зависимости от конфигурации версий и сервисов.
func (pr *ProxyRouter) RouterHandler(cfg *config.SocialNetworkConfig) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())

		requestMethod := c.Request().Method
		requestPath := c.Path()
		requestURL := c.Request().URL.String()
		acceptHeader := c.Request().Header.Get("Accept")
		version := parseProtocolVersion(acceptHeader)

		log.Debugw("Routing request",
			"requestMethod", requestMethod, "requestPath", requestPath, "requestURL", requestURL,
			"acceptHeader", acceptHeader, "version", version)

		// Поиск подходящего маршрута в конфиге
		for routeIdx, route := range cfg.RoutingConfig {
			if matchMethod(route.Methods, requestMethod) && matchRoute(route.Path, requestURL) {
				for serviceIdx, service := range route.Services {
					if isVersionSupported(service.SupportedVersions, version) {
						if proxy, exists := pr.proxies[route.Path][service.URL]; exists {
							// Здесь можно прописать circuit breaker
							log.Debugw("Route matched",
								"requestMethod", requestMethod, "requestPath", requestPath, "requestURL", requestURL,
								"acceptHeader", acceptHeader, "version", version, "routeIdx", routeIdx,
								"serviceIdx", serviceIdx, "service.ServiceName", service.ServiceName,
								"route.Path", route.Path, "version", version)
							return proxyRequest(c, proxy)
						}
					}
				}
			}
		}

		// Возвращаем 404, если маршрут не найден
		log.Warnw("No route matched",
			"requestMethod", requestMethod, "requestPath", requestPath, "requestURL", requestURL,
			"acceptHeader", acceptHeader, "version", version)
		return echo.NewHTTPError(http.StatusNotFound, "Route not found")
	}
}

// matchMethod проверяет, содержится ли requestMethod в списке pathMethods.
// Если pathMethods пуст, считаем, что метод допустим (разрешены все методы).
func matchMethod(pathMethods []string, requestMethod string) bool {
	if len(pathMethods) == 0 {
		return true
	}

	for _, method := range pathMethods {
		if method == requestMethod {
			return true
		}
	}
	return false
}

// parseProtocolVersion извлекает всё, что идет после 'application/' в заголовке Accept
func parseProtocolVersion(acceptHeader string) string {
	if acceptHeader == "" {
		return defaultProtocolVersion
	}

	if len(acceptHeader) > len(prefix) && acceptHeader[:len(prefix)] == prefix {
		return acceptHeader[len(prefix):]
	}

	// Если заголовок Accept не начинается с 'application/', возвращаем значение по умолчанию
	return defaultProtocolVersion
}

// matchRoute проверяет, соответствует ли путь маршруту (учитывая шаблоны)
func matchRoute(routePattern, requestPath string) bool {
	matched, err := regexp.MatchString(routePattern, requestPath)
	if err != nil {
		// В случае ошибки компиляции регулярного выражения считаем, что маршрут не совпадает
		return false
	}
	return matched
}

// isVersionSupported проверяет, поддерживает ли сервис указанную версию
func isVersionSupported(supportedVersions []string, version string) bool {
	for _, v := range supportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

// proxyRequest выполняет проксирование запроса в целевой сервис
func proxyRequest(c echo.Context, proxy *httputil.ReverseProxy) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())

	claims, ok := c.Get("claims").(*domain.UserClaims)
	if !ok {
		log.Warnw("Invalid claims in context")
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid claims in context")
	}

	userID := claims.Subject
	if userID == "" {
		log.Warnw("Missing user ID in claims")
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing user ID in claims")
	}

	hdr := c.Request().Header
	hdr.Set("X-User-Id", userID)
	// utils.AddRequestIDToHeader(ctx, c.Request().Header)

	// Здесь могут быть добавлены retry
	proxy.ServeHTTP(c.Response(), c.Request())

	return nil
}
