package rest

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
	"github.com/Vasiliy82/otus-hla-homework/internal/mappers"
	log "github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/internal/services"
	"github.com/Vasiliy82/otus-hla-homework/internal/validators"
	"github.com/labstack/echo/v4"
)

type postHandler struct {
	cfg         *config.PostHandlerConfig
	postService domain.PostService
}

func NewPostHandler(postService domain.PostService, cfg *config.PostHandlerConfig) services.PostHandler {
	return &postHandler{postService: postService, cfg: cfg}
}

// Создание нового поста (POST /posts)
func (h *postHandler) Create(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Create: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var postReq dto.CreateOrUpdatePostRequest

	if err = c.Bind(&postReq); err != nil {
		log.Logger().Errorw("postHandler.Create: c.Bind returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// Валидация запроса
	if err = validators.ValidateCreateOrUpdatePostRequest(postReq); err != nil {
		log.Logger().Warnw("postHandler.Create: validators.ValidateCreateOrUpdatePostRequest returned error", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	postMsg := mappers.ToPostMessage(&postReq)
	postId, err := h.postService.Create(userId, postMsg)
	if err != nil {
		log.Logger().Errorw("postHandler.Create: h.postService.Create returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.CreatePostResponse{Id: int64(postId)})
}

// Получение поста по ID (GET /posts/{id})
func (h *postHandler) Get(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Get: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	postId, err := getPostId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Get: getPostId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	post, err := h.postService.Get(userId, postId)
	if err != nil {
		log.Logger().Errorw("postHandler.Get: h.postService.Get returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mappers.ToPostResponse(post))
}

// Обновление поста по ID (PUT /posts/{id})
func (h *postHandler) Update(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Update: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var postReq dto.CreateOrUpdatePostRequest
	postId, err := getPostId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Update: getPostId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err := c.Bind(&postReq); err != nil {
		log.Logger().Errorw("postHandler.Update: c.Bind returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// Валидация запроса
	if err := validators.ValidateCreateOrUpdatePostRequest(postReq); err != nil {
		log.Logger().Warnw("postHandler.Update: validators.ValidateCreateOrUpdatePostRequest returned error", "err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.postService.Update(userId, postId, domain.PostMessage(postReq.Message))
	if err != nil {
		log.Logger().Errorw("postHandler.Update: h.postService.Update returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Удаление поста по ID (DELETE /posts/{id})
func (h *postHandler) Delete(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Delete: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	postId, err := getPostId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Delete: getPostId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err = h.postService.Delete(userId, postId)
	if err != nil {
		log.Logger().Errorw("postHandler.Delete: h.postService.Delete returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Получение ленты
func (h *postHandler) Feed(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		log.Logger().Errorw("postHandler.Get: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var limit int

	// Преобразование limit в int
	limitParam := c.QueryParam("limit")
	if limitParam == "" {
		limit = h.cfg.FeedDefaultPageSize
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 || limit > h.cfg.FeedMaxPageSize {
			log.Logger().Warnw("postHandler.Feed: invalid limit parameter", "limit", limitParam)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid limit parameter"})
		}
	}

	// Преобразование lastPostId в domain.PostKey
	var lastPostId domain.PostKey = math.MaxInt64

	lastPostIdParam := c.QueryParam("last_id")
	if lastPostIdParam != "" {
		lastPostIdInt, err := strconv.ParseInt(lastPostIdParam, 10, 64)
		if err != nil {
			log.Logger().Warnw("postHandler.Feed: invalid lastPostId parameter", "lastPostId", lastPostIdParam)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid lastPostId parameter"})
		}
		lastPostId = domain.PostKey(lastPostIdInt)
	}

	posts, err := h.postService.GetFeed(userId, limit, lastPostId)
	if err != nil {
		log.Logger().Errorw("postHandler.Feed: h.postService.Feed returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mappers.ToFeedResponse(posts))
}

func getUserId(c echo.Context) (domain.UserKey, error) {
	// Получаем информацию о пользователе
	claims, ok := c.Get("claims").(*domain.UserClaims)
	if !ok {
		// Теоретически, такого не должно случиться, т.к. токен проверяется в Middleware
		// log.Logger().Warnw("rest.getUserID: c.Get(\"claims\").(*domain.UserClaims) returned missing or invalid token")
		// залогируется выше
		return "", errors.New("missing or invalid token")
	}
	return domain.UserKey(claims.Subject), nil
}

func getPostId(c echo.Context) (domain.PostKey, error) {
	idParam := c.Param("post_id")
	postId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("rest.getPostId: strconv.ParseInt returned error: %w", err)
	}
	return domain.PostKey(postId), nil
}
