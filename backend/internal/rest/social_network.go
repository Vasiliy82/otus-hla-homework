package rest

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/dto"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/mappers"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/services"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/validators"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type socialNetworkHandler struct {
	snService domain.SocialNetworkService
	cfg       *config.APIConfig
}

// Регулярное выражение для проверки только букв
var validNameRegex = regexp.MustCompile(`^[\p{L}]+$`) // \p{L} соответствует любому юникодовскому символу, который является буквой

func NewSocialNetworkHandler(userService domain.SocialNetworkService, cfg *config.APIConfig) services.SocialNetworkHandler {
	return &socialNetworkHandler{
		snService: userService,
		cfg:       cfg,
	}
}

func (h *socialNetworkHandler) CreateUser(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	var userReq dto.RegisterUserRequest
	var user domain.User
	var err error
	log.Debug("UserHandler.RegisterUser")

	if err = c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err = validators.ValidateRegisterUserRequest(userReq); err != nil {
		var appverr *apperrors.ValidationError
		if errors.As(err, &appverr) {
			log.Warnw("userHandler.RegisterUser: validators.ValidateRegisterUserRequest returned ValidationError", "err", err)
			return c.JSON(http.StatusBadRequest, appverr)
		}
		log.Errorw("userHandler.RegisterUser: validators.ValidateRegisterUserRequest returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if user, err = mappers.ToUser(userReq); err != nil {
		log.Errorw("userHandler.RegisterUser: mappers.ToUser returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	userId, err := h.snService.CreateUser(&user)
	if err != nil {
		log.Errorw("userHandler.RegisterUser: h.userService.RegisterUser returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, apperr)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := dto.LoginResponse{ID: userId, Token: ""}

	return c.JSON(http.StatusOK, resp)
}

func (h *socialNetworkHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	var req dto.LoginRequest
	var err error

	log.Debug("Started")

	if err = c.Bind(&req); err != nil {
		log.Debug("c.Bind() returned error", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err = validators.ValidateLoginRequest(req); err != nil {
		log.Debug("validators.ValidateLoginRequest() returned error", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.snService.Login(ctx, req.Username, req.Password)

	if err != nil {
		log.Debug("h.snService.Login() returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	log.Debug("Finished")
	return c.JSON(http.StatusOK, map[string]string{"token": string(token)})
}

func (h *socialNetworkHandler) GetUser(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	var err error

	log.Debug("UserHandler.Get")

	id := domain.UserKey(c.Param("id"))

	if err = validators.ValidateUserId(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.snService.GetUser(id)
	if err != nil {
		log.Errorw("userHandler.Get: h.userService.GetById returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": apperr.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *socialNetworkHandler) Search(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	// Извлечение query параметров first_name и last_name
	firstName := c.QueryParam("first_name")
	lastName := c.QueryParam("last_name")

	// Валидация параметров
	if !isValidName(firstName) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат имени"})
	}
	if !isValidName(lastName) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неверный формат фамилии"})
	}

	users, err := h.snService.Search(firstName, lastName)
	if err != nil {
		log.Errorw("userHandler.Search: h.userService.Search returned error", "err", err)
		var apperr *apperrors.AppError
		if errors.As(err, &apperr) {
			return c.JSON(apperr.Code, map[string]string{"error": apperr.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *socialNetworkHandler) AddFriend(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	log.Debug("UserHandler.AddFriend")

	var err error

	// Извлечение идентификатора будущего друга из URL
	friend_id := domain.UserKey(c.Param("friend_id"))

	if err = validators.ValidateUserId(friend_id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Извлечение своего идентификатора из контекста
	// Здесь наверное нужен рефакторинг, сделать приватный метод со всеми проверками,
	// т.к. в будущем контекст пользователя понадобится практически везде
	claims, ok := c.Get("claims").(*domain.UserClaims)
	if !ok {
		// Теоретически, такого не должно случиться, т.к. токен проверяется в Middleware
		return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid token"))
	}

	my_id := domain.UserKey(claims.Subject)

	// Эту проверку имеет смысл вынести в middleware, чтобы валидировать в одном месте
	if err = validators.ValidateUserId(my_id); err != nil {
		log.Errorw("userHandler.AddFriend: validators.ValidateUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}

	if err = h.snService.AddFriend(my_id, friend_id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrFriendAlreadyExists) {
			log.Warnw("userHandler.AddFriend: h.userService.AddFriend returned ErrUserNotFound", "err", err)
			return c.JSON(http.StatusBadRequest, apperrors.NewBadRequestError(err.Error()))
		}
		log.Errorw("userHandler.AddFriend: h.userService.AddFriend returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *socialNetworkHandler) RemoveFriend(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	log.Debug("UserHandler.RemoveFriend")

	var err error

	// Извлечение идентификатора будущего друга из URL
	friend_id := domain.UserKey(c.Param("friend_id"))

	if err = validators.ValidateUserId(friend_id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Извлечение своего идентификатора из контекста
	// Здесь наверное нужен рефакторинг, сделать приватный метод со всеми проверками,
	// т.к. в будущем контекст пользователя понадобится практически везде
	claims, ok := c.Get("claims").(*domain.UserClaims)
	if !ok {
		// Теоретически, такого не должно случиться, т.к. токен проверяется в Middleware
		return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid token"))
	}

	my_id := domain.UserKey(claims.Subject)

	// Эту проверку имеет смысл вынести в middleware, чтобы валидировать в одном месте
	if err = validators.ValidateUserId(my_id); err != nil {
		log.Errorw("userHandler.RemoveFriend: validators.ValidateUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}

	if err = h.snService.RemoveFriend(my_id, friend_id); err != nil {
		if errors.Is(err, domain.ErrFriendNotFound) {
			log.Warnw("userHandler.RemoveFriend: h.userService.RemoveFriend returned ErrFriendNotFound", "err", err)
			return c.JSON(http.StatusBadRequest, apperrors.NewBadRequestError(err.Error()))
		}
		log.Errorw("userHandler.RemoveFriend: h.userService.RemoveFriend returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *socialNetworkHandler) Logout(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	log.Debug("UserHandler.Logout")
	// Извлекаем токен из контекста
	token, ok := c.Get("token").(*jwt.Token)
	if !ok {
		// Теоретически такого не может случиться, т.к. токен проверяется в Middleware
		return c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorizedError("missing or invalid token"))
	}

	if err := h.snService.Logout(token); err != nil {
		log.Errorw("userHandler.Logout: h.userService.Logout returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, apperrors.NewInternalServerError("Internal server error", err))
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *socialNetworkHandler) ListPosts(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	userId, err := getUserId(c)
	if err != nil {
		log.Errorw("postHandler.List: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	limit, lastPostId, err := h.getLimits(c)
	if err != nil {
		log.Errorw("postHandler.Feed: h.getLimits returned error", "err", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	posts, err := h.snService.ListPosts(userId, limit, lastPostId)
	if err != nil {
		log.Errorw("postHandler.Feed: h.snService.ListPosts returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, mappers.ToPostsResponse(posts))
}

// Создание нового поста (POST /posts)
func (h *socialNetworkHandler) CreatePost(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	userId, err := getUserId(c)
	if err != nil {
		log.Errorw("postHandler.Create: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var postReq dto.CreateOrUpdatePostRequest

	if err = c.Bind(&postReq); err != nil {
		log.Errorw("postHandler.Create: c.Bind returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// Валидация запроса
	if err = validators.ValidateCreateOrUpdatePostRequest(postReq); err != nil {
		log.Warnw("postHandler.Create: validators.ValidateCreateOrUpdatePostRequest returned error", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	postMsg := mappers.ToPostMessage(&postReq)
	postId, err := h.snService.CreatePost(userId, postMsg)
	if err != nil {
		log.Errorw("postHandler.Create: h.postService.Create returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.CreatePostResponse{Id: int64(postId)})
}

// Получение поста по ID (GET /posts/{id})
func (h *socialNetworkHandler) GetPost(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	userId, err := getUserId(c)
	if err != nil {
		log.Errorw("postHandler.Get: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	postId, err := getPostId(c)
	if err != nil {
		log.Errorw("postHandler.Get: getPostId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	post, err := h.snService.GetPost(userId, postId)
	if err != nil {
		log.Errorw("postHandler.Get: h.postService.Get returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, mappers.ToPostResponse(post))
}

// Обновление поста по ID (PUT /posts/{id})
func (h *socialNetworkHandler) UpdatePost(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	userId, err := getUserId(c)
	if err != nil {
		log.Errorw("postHandler.Update: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var postReq dto.CreateOrUpdatePostRequest
	postId, err := getPostId(c)
	if err != nil {
		log.Errorw("postHandler.Update: getPostId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if err := c.Bind(&postReq); err != nil {
		log.Errorw("postHandler.Update: c.Bind returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// Валидация запроса
	if err := validators.ValidateCreateOrUpdatePostRequest(postReq); err != nil {
		log.Warnw("postHandler.Update: validators.ValidateCreateOrUpdatePostRequest returned error", "err", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.snService.UpdatePost(userId, postId, domain.PostText(postReq.Message))
	if err != nil {
		log.Errorw("postHandler.Update: h.postService.Update returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Удаление поста по ID (DELETE /posts/{id})
func (h *socialNetworkHandler) DeletePost(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	userId, err := getUserId(c)
	if err != nil {
		log.Errorw("postHandler.Delete: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	postId, err := getPostId(c)
	if err != nil {
		log.Errorw("postHandler.Delete: getPostId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err = h.snService.DeletePost(userId, postId)
	if err != nil {
		log.Errorw("postHandler.Delete: h.postService.Delete returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Получение ленты
func (h *socialNetworkHandler) GetFeed(c echo.Context) error {
	log := logger.FromContext(c.Request().Context()).With("func", logger.GetFuncName())
	userId, err := getUserId(c)
	if err != nil {
		log.Errorw("postHandler.Feed: getUserId returned error", "err", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	posts, err := h.snService.GetFeed(userId)
	if err != nil {
		log.Errorw("postHandler.Feed: h.snService.GetFeed returned error", "err", err)
		if errors.Is(err, domain.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, mappers.ToPostsResponse(posts))
}

func getUserId(c echo.Context) (domain.UserKey, error) {
	// Получаем информацию о пользователе
	claims, ok := c.Get("claims").(*domain.UserClaims)
	if !ok {
		// Теоретически, такого не должно случиться, т.к. токен проверяется в Middleware
		// log.Warnw("rest.getUserID: c.Get(\"claims\").(*domain.UserClaims) returned missing or invalid token")
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

func (h socialNetworkHandler) getLimits(c echo.Context) (int, domain.PostKey, error) {
	var limit int
	var err error
	var lastPostId domain.PostKey = math.MaxInt64

	// Преобразование limit в int
	limitParam := c.QueryParam("limit")
	if limitParam == "" {
		limit = h.cfg.FeedDefaultPageSize
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 || limit > h.cfg.FeedMaxPageSize {
			return 0, 0, fmt.Errorf("invalid limit parameter")
		}
	}

	lastPostIdParam := c.QueryParam("last_id")
	if lastPostIdParam != "" {
		lastPostIdInt, err := strconv.ParseInt(lastPostIdParam, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid lastPostId parameter")
		}
		lastPostId = domain.PostKey(lastPostIdInt)
	}
	return limit, lastPostId, nil
}

// Функция для валидации имени
func isValidName(name string) bool {
	// Проверяем, что строка содержит только буквы
	return validNameRegex.MatchString(name)
}
