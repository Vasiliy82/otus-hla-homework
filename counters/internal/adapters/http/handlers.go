package http

import (
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
	"github.com/labstack/echo/v4"
)

type CounterHTTPHandler struct {
	counterService domain.CounterService
}

func NewCounterHTTPHandler(counterService domain.CounterService) *CounterHTTPHandler {
	return &CounterHTTPHandler{counterService: counterService}
}

func (h *CounterHTTPHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/counters/:dialogID/increment", h.IncrementUnread)
	e.POST("/counters/:dialogID/reset", h.ResetUnread)
}

func (h *CounterHTTPHandler) IncrementUnread(c echo.Context) error {
	dialogID := c.Param("dialogID")
	if err := h.counterService.IncrementUnread(c.Request().Context(), dialogID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "incremented"})
}

func (h *CounterHTTPHandler) ResetUnread(c echo.Context) error {
	dialogID := c.Param("dialogID")
	if err := h.counterService.ResetUnread(c.Request().Context(), dialogID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "reset"})
}
