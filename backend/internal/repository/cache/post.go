package cache

import (
	"context"
	"encoding/json"

	"github.com/Vasiliy82/otus-hla-homework/backend/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

type postCache struct {
	client *redis.Client
	cfg    *config.CacheConfig
}

// NewPostCache создает новый экземпляр postCache
func NewPostCache(cfg *config.CacheConfig, client *redis.Client) domain.PostCache {
	return &postCache{
		client: client,
		cfg:    cfg,
	}
}

// GetFeed возвращает ленту пользователя из кеша, декодируя JSON-строку в массив постов
func (c *postCache) GetFeed(userId domain.UserKey, limit int) ([]*domain.Post, error) {
	key := getCacheKey(userId)

	// Получаем JSON данные из Redis
	data, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// Ключ отсутствует, данные нужно получить из SQL и добавить в кеш
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Если данные присутствуют, декодируем JSON
	var posts []*domain.Post
	if err := json.Unmarshal([]byte(data), &posts); err != nil {
		return nil, err
	}

	// Ограничиваем количество возвращаемых постов
	if len(posts) > limit {
		posts = posts[:limit]
	}

	if posts == nil {
		posts = []*domain.Post{}
	}

	return posts, nil
}

// UpdateFeed сохраняет ленту пользователя в кеше в формате JSON
func (c *postCache) UpdateFeed(userId domain.UserKey, posts []*domain.Post) error {
	key := getCacheKey(userId)

	// Преобразуем массив постов в JSON
	jsonData, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	// Сохраняем JSON строку в Redis с указанным временем жизни
	return c.client.Set(context.Background(), key, jsonData, c.cfg.Expiry).Err()
}

// getCacheKey формирует ключ кеша для пользователя
func getCacheKey(userId domain.UserKey) string {
	return "feed:" + string(userId)
}
