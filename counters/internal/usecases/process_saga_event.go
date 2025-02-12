package usecases

import (
	"context"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
)

// ProcessSagaEventUseCase управляет обработкой событий SAGA
type ProcessSagaEventUseCase struct {
	counterService domain.CounterService
	publisher      domain.KafkaPublisher
}

func NewProcessSagaEventUseCase(counterService domain.CounterService, publisher domain.KafkaPublisher) *ProcessSagaEventUseCase {
	return &ProcessSagaEventUseCase{
		counterService: counterService,
		publisher:      publisher,
	}
}

// Execute обрабатывает события SAGA
func (uc *ProcessSagaEventUseCase) Execute(ctx context.Context, event domain.SagaEvent) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "event", event)

	switch event.Type {
	case domain.SagaMessageSent:
		err := uc.handleMessageSent(ctx, event)
		if err != nil {
			log.Warnw("Failure", "event", event, "err", err)
			return err
		}
		log.Infow("Success", "event", event)
		return nil
	case domain.SagaCounterIncremented:
		log.Debugw("Skipped")
		return nil
	case domain.SagaCounterIncrementFailed:
		log.Debugw("Skipped")
		return nil
	default:
		err := fmt.Errorf("unknown SAGA event type: %s", event.Type)
		log.Warnw("Failure", "event", event, "err", err)
		return err
	}
}

// handleMessageSent обрабатывает событие отправки сообщения и увеличивает счетчик
func (uc *ProcessSagaEventUseCase) handleMessageSent(ctx context.Context, event domain.SagaEvent) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("Started", "event", event)

	// Пытаемся увеличить счетчик
	err := uc.counterService.IncrementUnread(ctx, event.DialogID)
	if err != nil {
		log.Debugw("uc.counterService.IncrementUnread() returned error", "err", err)
		// Если ошибка, отправляем событие о неудаче
		failEvent := domain.SagaEvent{
			TransactionID: event.TransactionID,
			Type:          "CounterIncrementFailed",
			DialogID:      event.DialogID,
			Error:         err.Error(),
		}
		err1 := uc.publisher.PublishSagaEvent(ctx, failEvent)
		log.Debugw("uc.publisher.PublishSagaEvent()", "failEvent", failEvent, "err", err1)
		return err
	}

	// Если успех, отправляем подтверждающее событие
	successEvent := domain.SagaEvent{
		TransactionID: event.TransactionID,
		Type:          "CounterIncremented",
		DialogID:      event.DialogID,
	}
	err = uc.publisher.PublishSagaEvent(ctx, successEvent)
	log.Debugw("uc.publisher.PublishSagaEvent()", "successEvent", successEvent, "err", err)
	return err
}
