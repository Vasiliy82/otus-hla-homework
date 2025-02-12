package repository

import (
	"context"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
)

type CombinedCounterRepository struct {
	redisRepo *RedisCounterRepository
	pgRepo    *PGCounterRepository
}

func NewCombinedCounterRepository(redisRepo *RedisCounterRepository, pgRepo *PGCounterRepository) *CombinedCounterRepository {
	return &CombinedCounterRepository{
		redisRepo: redisRepo,
		pgRepo:    pgRepo,
	}
}

func (r *CombinedCounterRepository) GetCounter(ctx context.Context, dialogID string) (*domain.DialogCounter, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "dialogID", dialogID)
	// Попытка получить данные из Redis
	counter, err := r.redisRepo.GetCounter(ctx, dialogID)
	if err == nil && counter != nil {
		log.Debugw("Success", "counter", counter)
		return counter, nil
	}

	// Если в Redis нет данных, обращаемся к PostgreSQL
	counter, err = r.pgRepo.GetCounter(ctx, dialogID)
	if err != nil {
		log.Debugw("r.pgRepo.GetCounter() returned error", "err", err)
		return nil, err
	}

	// После получения из PostgreSQL, кешируем результат в Redis
	_ = r.redisRepo.SaveCounter(ctx, counter)

	log.Debugw("Success", "counter", counter)
	return counter, nil
}

func (r *CombinedCounterRepository) SaveCounter(ctx context.Context, counter *domain.DialogCounter) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	// Сохраняем данные в Redis и PostgreSQL
	if err := r.redisRepo.SaveCounter(ctx, counter); err != nil {
		log.Debug("r.redisRepo.SaveCounter() returned error", "err", err)
		return err
	}
	log.Debugw("Success", "counter", counter)
	return r.pgRepo.SaveCounter(ctx, counter)
}
