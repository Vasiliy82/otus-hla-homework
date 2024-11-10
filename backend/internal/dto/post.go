package dto

import (
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
)

// Структура запроса для создания нового поста
type CreateOrUpdatePostRequest struct {
	Message string `json:"message" validate:"required"` // Сообщение обязательно
}

// Структура ответа для получения поста
type GetPostResponse struct {
	Id         int64      `json:"id"`
	Message    string     `json:"message"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"` // Время последнего редактирования может отсутствовать
}

// Структура ответа при создании поста
type CreatePostResponse struct {
	Id int64 `json:"id"` // ID поста
}

// Структура ответа для получения ленты новостей
type GetFeedResponse struct {
	Feed       []*domain.Post `json:"feed"`
	LastPostId domain.PostKey `json:"last_id"`
}
