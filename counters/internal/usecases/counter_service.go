package usecases

import (
	"context"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
)

// CounterService реализует интерфейс domain.CounterService
type CounterService struct {
	repo domain.CounterRepository
}

// Конструктор для CounterService
func NewCounterService(repo domain.CounterRepository) *CounterService {
	return &CounterService{repo: repo}
}

// IncrementUnread увеличивает счетчик непрочитанных сообщений
func (cs *CounterService) IncrementUnread(ctx context.Context, dialogID string) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "dialogID", dialogID)

	counter, err := cs.repo.GetCounter(ctx, dialogID)
	if err != nil {
		log.Warnw("Failure", "dialogID", dialogID, "err", err)
		return err
	}

	counter.UnreadCount++

	if counter.UnreadCount < 0 {
		log.Warnw("Failure", "dialogID", dialogID, "err", domain.ErrCounterOverflow)
		return domain.ErrCounterOverflow
	}
	log.Infow("Success", "counter", counter)
	return cs.repo.SaveCounter(ctx, counter)
}

// ResetUnread сбрасывает счетчик непрочитанных сообщений в 0
func (cs *CounterService) ResetUnread(ctx context.Context, dialogID string) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "dialogID", dialogID)

	counter, err := cs.repo.GetCounter(ctx, dialogID)
	if err != nil {
		log.Warnw("Failure", "dialogID", dialogID, "err", err)
		return err
	}
	counter.UnreadCount = 0
	err = cs.repo.SaveCounter(ctx, counter)
	if err != nil {
		log.Warnw("Failure", "dialogID", dialogID, "err", err)
		return err
	}
	log.Infow("Success", "counter", counter)
	return nil
}
