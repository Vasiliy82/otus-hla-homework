package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/internal/infrastructure/broker"
	log "github.com/Vasiliy82/otus-hla-homework/internal/observability/logger"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type SocialNetworkHandler interface {
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	Search(c echo.Context) error
	Login(c echo.Context) error
	AddFriend(c echo.Context) error
	RemoveFriend(c echo.Context) error
	Logout(c echo.Context) error
	ListPosts(c echo.Context) error
	CreatePost(c echo.Context) error
	GetPost(c echo.Context) error
	UpdatePost(c echo.Context) error
	DeletePost(c echo.Context) error
	GetFeed(c echo.Context) error
}

type socialNetworkService struct {
	userRepo   domain.UserRepository
	postRepo   domain.PostRepository
	postCache  domain.PostCache
	jwtService domain.JWTService
	cfg        *config.SocialNetworkConfig
	producer   *broker.Producer
}

func NewSocialNetworkService(cfg *config.Config,
	ur domain.UserRepository,
	pr domain.PostRepository,
	pc domain.PostCache,
	jwts domain.JWTService,
	producer *broker.Producer) domain.SocialNetworkService {
	return &socialNetworkService{
		userRepo:   ur,
		postRepo:   pr,
		postCache:  pc,
		jwtService: jwts,
		cfg:        cfg.SocialNetwork,
		producer:   producer,
	}
}

func (s *socialNetworkService) CreateUser(user *domain.User) (domain.UserKey, error) {
	var id domain.UserKey
	var err error

	if id, err = s.userRepo.RegisterUser(user); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // duplicate key value violates unique constraint
				return "", apperrors.NewConflictError("Login already used")
			}
		}
		// Если ошибка не является *pq.Error, оборачиваем её в InternalServerError
		return "", apperrors.NewInternalServerError("UserService.RegisterUser, s.userRepo.RegisterUser returned unknown error", err)
	}
	return id, nil
}

func (s *socialNetworkService) GetUser(id domain.UserKey) (*domain.User, error) {
	var user *domain.User
	var err error
	if user, err = s.userRepo.GetByID(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.User{}, apperrors.NewNotFoundError("User not found")
		}
		return &domain.User{}, apperrors.NewInternalServerError("UserService.GetById: s.userRepo.GetByID returned unknown error", err)
	}
	return user, nil

}

func (s *socialNetworkService) Login(username, password string) (domain.TokenString, error) {
	// Проверка пароля
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperrors.NewNotFoundError("User not found")
		}
		return "", apperrors.NewInternalServerError("UserService.Login: s.userRepo.GetByUserName returned unknown error", err)
	}
	if !user.CheckPassword(password) {
		return "", apperrors.NewUnauthorizedError("Wrong password")
	}

	token, err := s.jwtService.GenerateToken(user.ID, []domain.Permission{domain.PermissionUserGet})
	if err != nil {
		return "", apperrors.NewInternalServerError("UserSevice.Login: s.sessionRepo.CreateSession returned unknown error", err)
	}

	return token, nil
}

func (s *socialNetworkService) Search(firstName, lastName string) ([]*domain.User, error) {
	users, err := s.userRepo.Search(firstName, lastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFoundError("User not found")
		}
		return nil, apperrors.NewInternalServerError("UserService.Login: s.userRepo.GetByUserName returned unknown error", err)
	}

	return users, nil
}

func (s *socialNetworkService) AddFriend(my_id, friend_id domain.UserKey) error {

	if err := s.userRepo.AddFriend(my_id, friend_id); err != nil {
		if err == domain.ErrObjectAlreadyExists {
			return domain.ErrFriendAlreadyExists
		}
		if err == domain.ErrObjectNotFound {
			return domain.ErrUserNotFound
		}
		return apperrors.NewInternalServerError("Internal server error", err)
	}
	return nil
}

func (s *socialNetworkService) RemoveFriend(my_id, friend_id domain.UserKey) error {

	if err := s.userRepo.RemoveFriend(my_id, friend_id); err != nil {
		if err == domain.ErrObjectNotFound {
			return domain.ErrFriendNotFound
		}
		return apperrors.NewInternalServerError("Internal server error", err)
	}
	return nil
}

func (s *socialNetworkService) Logout(token *jwt.Token) error {

	if err := s.jwtService.RevokeToken(token); err != nil {
		return apperrors.NewInternalServerError("Internal server error", err)
	}
	return nil
}

func (s *socialNetworkService) ListPosts(userId domain.UserKey, limit int, lastPostId domain.PostKey) ([]*domain.Post, error) {
	posts, err := s.postRepo.List(userId, limit, lastPostId)
	if err != nil {
		return nil, apperrors.NewInternalServerError("postService.List: s.postRepo.List returned error", err)
	}
	return posts, nil
}

// Создание нового поста
func (s *socialNetworkService) CreatePost(userId domain.UserKey, message domain.PostMessage) (domain.PostKey, error) {
	postId, err := s.postRepo.Create(userId, message)
	if err != nil {
		return 0, apperrors.NewInternalServerError("postService.Create: s.postRepo.Create returned error", err)
	}
	s.sendRecalculationEvent(userId, domain.EventPostCreated)
	return postId, nil
}

// Получение поста по ID
func (s *socialNetworkService) GetPost(userId domain.UserKey, id domain.PostKey) (*domain.Post, error) {
	post, err := s.postRepo.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFoundError("Post not found")
		}
		return nil, apperrors.NewInternalServerError("postService.Get: s.postRepo.GetByID returned error", err)
	}
	return post, nil
}

// Обновление поста
func (s *socialNetworkService) UpdatePost(userId domain.UserKey, postId domain.PostKey, newMessage domain.PostMessage) error {
	if err := s.checkOwner(userId, postId); err != nil {
		return apperrors.New(403, "Wrong post owner", nil)
	}

	err := s.postRepo.UpdateMessage(postId, newMessage)
	if err != nil {
		if errors.Is(err, domain.ErrObjectNotFound) {
			return apperrors.NewNotFoundError("Post not found")
		}
		return apperrors.NewInternalServerError("postService.Update: s.postRepo.Update returned error", err)
	}
	s.sendRecalculationEvent(userId, domain.EventPostEdited)
	return nil
}

// Удаление поста
func (s *socialNetworkService) DeletePost(userId domain.UserKey, postId domain.PostKey) error {
	if err := s.checkOwner(userId, postId); err != nil {
		return apperrors.New(403, "Wrong post owner", nil)
	}

	err := s.postRepo.Delete(postId)
	if err != nil {
		if errors.Is(err, domain.ErrObjectNotFound) {
			return apperrors.NewNotFoundError("Post not found")
		}
		return apperrors.NewInternalServerError("postService.Delete: s.postRepo.Delete returned error", err)
	}
	s.sendRecalculationEvent(userId, domain.EventPostDeleted)
	return nil
}

func (s *socialNetworkService) GetFeed(userId domain.UserKey) ([]*domain.Post, error) {
	// Если есть кеш, то сначала посмотрим там
	if s.postCache != nil {
		cache, err := s.postCache.GetFeed(userId, s.cfg.FeedLength)
		if err != nil {
			log.Logger().Warnw("postService.Feed: s.postCache.GetFeed returned error", "err", err)
		} else if cache != nil {
			return cache, nil
		}
	}

	// Если в кеше ничего нет (забыли или не успели "прогреть", бывает), то берем из БД
	posts, err := s.postRepo.GetFeed(userId, s.cfg.FeedLength)
	if err != nil {
		return nil, apperrors.NewInternalServerError("postService.Feed: s.postRepo.GetFeed returned error", err)
	}
	if s.postCache != nil {
		if err = s.postCache.UpdateFeed(userId, posts); err != nil {
			log.Logger().Warnw("postService.Feed: s.postCache.UpdateFeed returned error", "err", err)
		}
	}
	return posts, nil
}

func (s *socialNetworkService) SetLastActivity(userId domain.UserKey) error {
	if err := s.userRepo.SetLastActivity(userId); err != nil {
		return fmt.Errorf("socialNetworkService.SetLastActivity: s.userRepo.SetLastActivity() returned error: %w", err)
	}
	return nil
}

func (s *socialNetworkService) checkOwner(userId domain.UserKey, postId domain.PostKey) error {
	postOwner, err := s.postRepo.GetPostOwner(postId)
	if err != nil {
		return apperrors.NewInternalServerError("postService.checkOwner: s.postRepo.GetPostOwner returned error", err)
	}
	if postOwner != userId {
		return apperrors.New(403, "Wrong post owner", nil)
	}
	return nil
}

func (s *socialNetworkService) sendRecalculationEvent(userId domain.UserKey, eventType domain.EventType) error {
	if s.producer == nil {
		return nil
	}

	// Создаем и сериализуем событие
	event := domain.EventInvalidateCache{
		UserID:    userId,
		EventType: eventType,
	}

	// Отправляем событие в Kafka
	if err := s.producer.SendCacheEvent(event); err != nil {
		log.Logger().Errorw("Ошибка отправки события в Kafka", "userId", userId, "err", err)
		return err
	}

	return nil
}
