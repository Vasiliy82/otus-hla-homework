package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisCounterRepository struct {
	client *redis.Client
}

func NewRedisCounterRepository(addr string, port int, password string, db int) *RedisCounterRepository {
	addrport := fmt.Sprintf("%s:%d", addr, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addrport,
		Password: password,
		DB:       db,
	})
	return &RedisCounterRepository{client: client}
}

func (r *RedisCounterRepository) GetCounter(ctx context.Context, dialogID string) (*domain.DialogCounter, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "dialogID", dialogID)
	data, err := r.client.Get(ctx, dialogID).Result()
	// if err == redis.Nil {
	// 	return &domain.DialogCounter{DialogID: dialogID, UnreadCount: 0}, nil
	// } else
	if err != nil {
		log.Debugw("r.client.Get() returned error", "err", err)
		return nil, err
	}

	var counter domain.DialogCounter
	if err := json.Unmarshal([]byte(data), &counter); err != nil {
		log.Debugw("json.Unmarshal() returned error", "err", err)
		return nil, err
	}
	log.Debugw("Success", "counter", counter)
	return &counter, nil
}

func (r *RedisCounterRepository) SaveCounter(ctx context.Context, counter *domain.DialogCounter) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "counter", counter)
	data, err := json.Marshal(counter)
	if err != nil {
		log.Debugw("json.Marshal() returned error", "err", err)
		return err
	}

	err = r.client.Set(ctx, counter.DialogID, data, 0).Err()
	if err != nil {
		log.Debugw(".client.Set().Err() returned error", "err", err)
		return err
	}
	log.Debug("Success")
	return nil
}
