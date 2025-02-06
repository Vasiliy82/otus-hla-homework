package rest

import (
	"errors"
	"net/http"
	"strconv"

	apperrors "github.com/Vasiliy82/otus-hla-homework/backend/internal/apperrors2"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/labstack/echo/v4"
)

type dialogHandler struct {
	dialogService domain.DialogService
}

// NewDialogHandler создает новый экземпляр обработчика
func NewDialogHandler(dialogService domain.DialogService) *dialogHandler {
	return &dialogHandler{dialogService: dialogService}
}

// SendMessage обрабатывает POST /dialog/:partnerId/send
func (h *dialogHandler) SendMessage(c echo.Context) error {
	partnerId := domain.UserKey(c.Param("partnerId"))

	var req struct {
		Text string `json:"text"`
	}
	if err := c.Bind(&req); err != nil {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid request body",
			Err:     err,
		})
	}

	// Извлечение user_id из заголовка X-User-Id
	userIdHeader := c.Request().Header.Get("X-User-Id")
	if userIdHeader == "" {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ClientError,
			Message: "missing X-User-Id header",
		})
	}

	myId := domain.UserKey(userIdHeader)

	err := h.dialogService.SendMessage(c.Request().Context(), myId, partnerId, req.Text)
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetDialog обрабатывает GET /dialog/:partnerId/list
func (h *dialogHandler) GetDialog(c echo.Context) error {
	partnerId := domain.UserKey(c.Param("partnerId"))

	// Извлечение user_id из заголовка X-User-Id
	userIdHeader := c.Request().Header.Get("X-User-Id")
	if userIdHeader == "" {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ClientError,
			Message: "missing X-User-Id header",
		})
	}

	myId := domain.UserKey(userIdHeader)

	// Чтение offset и limit из query parameters
	offset, err := parseQueryParam(c, "offset", 0)
	if err != nil {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid offset value",
			Err:     err,
		})
	}

	limit, err := parseQueryParam(c, "limit", 20) // Значение по умолчанию — 20
	if err != nil {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid limit value",
			Err:     err,
		})
	}

	dialog, err := h.dialogService.GetDialog(c.Request().Context(), myId, partnerId, limit, offset)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dialog)
}

func (h *dialogHandler) GetDialogs(c echo.Context) error {

	// Извлечение user_id из заголовка X-User-Id
	userIdHeader := c.Request().Header.Get("X-User-Id")
	if userIdHeader == "" {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ClientError,
			Message: "missing X-User-Id header",
		})
	}

	myId := domain.UserKey(userIdHeader)

	// Чтение offset и limit из query parameters
	offset, err := parseQueryParam(c, "offset", 0)
	if err != nil {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid offset value",
			Err:     err,
		})
	}

	limit, err := parseQueryParam(c, "limit", 20) // Значение по умолчанию — 20
	if err != nil {
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid limit value",
			Err:     err,
		})
	}

	dialogs, err := h.dialogService.GetDialogs(c.Request().Context(), myId, limit, offset)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dialogs)
}

// parseQueryParam парсит числовые параметры из query string
func parseQueryParam(c echo.Context, key string, defaultValue int) (int, error) {
	param := c.QueryParam(key)
	if param == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	if value < 0 {
		return 0, errors.New("value cannot be negative")
	}

	return value, nil
}

// handleError обрабатывает ошибку, логирует цепочку и возвращает корректный HTTP-ответ
func handleError(c echo.Context, err error) error {
	// Логируем всю цепочку ошибок
	logErrorChain(c, err)

	// Пытаемся найти бизнес-ошибку (AppError)
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		// Если это известная бизнес-ошибка, возвращаем ее сообщение
		switch appErr.Type {
		case apperrors.ClientError, apperrors.ValidationError:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": appErr.Message})
		case apperrors.AuthorizationError:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": appErr.Message})
		case apperrors.RemoteServiceError:
			return c.JSON(http.StatusBadGateway, map[string]string{"error": appErr.Message})
		}
	}

	// Если бизнес-ошибка не найдена, возвращаем общую ошибку сервиса
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

// logErrorChain логирует всю цепочку ошибок
func logErrorChain(c echo.Context, err error) {
	for err != nil {
		c.Logger().Error(err.Error()) // Используем логгер Echo
		err = errors.Unwrap(err)
	}
}
