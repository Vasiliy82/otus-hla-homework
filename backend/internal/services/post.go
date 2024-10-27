package services

import (
	"database/sql"
	"errors"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/labstack/echo/v4"
)

type postService struct {
	postRepo domain.PostRepository
}

type PostHandler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Feed(c echo.Context) error
}

func NewPostService(ur domain.PostRepository) domain.PostService {
	return &postService{
		postRepo: ur,
	}
}

// Создание нового поста
func (s *postService) Create(userId domain.UserKey, message domain.PostMessage) (domain.PostKey, error) {
	postId, err := s.postRepo.Create(userId, message)
	if err != nil {
		return 0, apperrors.NewInternalServerError("postService.Create: s.postRepo.Create returned error", err)
	}
	return postId, nil
}

// Получение поста по ID
func (s *postService) Get(userId domain.UserKey, id domain.PostKey) (*domain.Post, error) {
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
func (s *postService) Update(userId domain.UserKey, postId domain.PostKey, newMessage domain.PostMessage) error {
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
	return nil
}

// Удаление поста
func (s *postService) Delete(userId domain.UserKey, postId domain.PostKey) error {
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
	return nil
}

func (s *postService) GetFeed(userId domain.UserKey, limit int, lastPostId domain.PostKey) ([]*domain.Post, error) {
	posts, err := s.postRepo.GetFeed(userId, limit, lastPostId)
	if err != nil {
		return nil, apperrors.NewInternalServerError("postService.Feed: s.postRepo.Feed returned error", err)
	}
	return posts, nil
}

func (s *postService) checkOwner(userId domain.UserKey, postId domain.PostKey) error {
	postOwner, err := s.postRepo.GetPostOwner(postId)
	if err != nil {
		return apperrors.NewInternalServerError("postService.checkOwner: s.postRepo.GetPostOwner returned error", err)
	}
	if postOwner != userId {
		return apperrors.New(403, "Wrong post owner", nil)
	}
	return nil
}
