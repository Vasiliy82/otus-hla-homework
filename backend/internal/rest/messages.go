package rest

import (
	"errors"
	"net/http"
	"strconv"

	apperrors "github.com/Vasiliy82/otus-hla-homework/backend/internal/apperrors2"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
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
	req := c.Request()
	ctx := req.Context()
	log := logger.FromContext(ctx)

	partnerId := domain.UserKey(c.Param("partnerId"))

	var reqBody struct {
		Text string `json:"text"`
	}
	if err := c.Bind(&reqBody); err != nil {
		log.Warnw("Error parsing request body", "reqBody", reqBody, "err", err)
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid request body",
			Err:     err,
		})
	}

	// Извлечение user_id из заголовка X-User-Id
	userIdHeader := c.Request().Header.Get("X-User-Id")
	if userIdHeader == "" {
		log.Warn("missing X-User-Id header")
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ClientError,
			Message: "missing X-User-Id header",
		})
	}

	myId := domain.UserKey(userIdHeader)
	ctx = logger.WithContext(ctx, logger.FromContext(ctx).With("userID", myId))

	err := h.dialogService.SendMessage(ctx, myId, partnerId, reqBody.Text)
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetDialog обрабатывает GET /dialog/:partnerId/list
func (h *dialogHandler) GetDialog(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	partnerId := domain.UserKey(c.Param("partnerId"))

	// Извлечение user_id из заголовка X-User-Id
	userIdHeader := req.Header.Get("X-User-Id")
	if userIdHeader == "" {
		log.Warn("missing X-User-Id header")
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ClientError,
			Message: "missing X-User-Id header",
		})
	}

	myId := domain.UserKey(userIdHeader)
	ctx = logger.WithContext(ctx, logger.FromContext(ctx).With("userID", myId))

	// Чтение offset и limit из query parameters
	offset, err := parseQueryParam(c, "offset", 0)
	if err != nil {
		log.Warn("invalid offset value")
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid offset value",
			Err:     err,
		})
	}

	limit, err := parseQueryParam(c, "limit", 20) // Значение по умолчанию — 20
	if err != nil {
		log.Warn("invalid limit value")
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid limit value",
			Err:     err,
		})
	}

	dialog, err := h.dialogService.GetDialog(ctx, myId, partnerId, limit, offset)
	if err != nil {
		log.Warn("h.dialogService.GetDialog() returned error")
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, dialog)
}

func (h *dialogHandler) GetDialogs(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())

	// Извлечение user_id из заголовка X-User-Id
	userIdHeader := req.Header.Get("X-User-Id")
	if userIdHeader == "" {
		log.Warn("missing X-User-Id header")
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ClientError,
			Message: "missing X-User-Id header",
		})
	}

	myId := domain.UserKey(userIdHeader)
	ctx = logger.WithContext(ctx, logger.FromContext(ctx).With("userID", myId))

	// Чтение offset и limit из query parameters
	offset, err := parseQueryParam(c, "offset", 0)
	if err != nil {
		log.Warn("invalid offset value")
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid offset value",
			Err:     err,
		})
	}

	limit, err := parseQueryParam(c, "limit", 20) // Значение по умолчанию — 20
	if err != nil {
		log.Warn("invalid limit value")
		return handleError(c, &apperrors.AppError{
			Type:    apperrors.ValidationError,
			Message: "invalid limit value",
			Err:     err,
		})
	}

	dialogs, err := h.dialogService.GetDialogs(ctx, myId, limit, offset)
	if err != nil {
		log.Warn("h.dialogService.GetDialogs() returned error", "err", err)
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
