package usecases

import (
	"context"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/broker"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
)

// SagaOrchestrator управляет процессами SAGA и реализует SagaCoordinator
type SagaOrchestrator struct {
	kafkaProducer *broker.Producer
	dialogService domain.DialogService
}

// NewSagaOrchestrator создает новый SagaOrchestrator
func NewSagaOrchestrator(kafkaProducer *broker.Producer, dialogService domain.DialogService) *SagaOrchestrator {
	return &SagaOrchestrator{
		kafkaProducer: kafkaProducer,
		dialogService: dialogService,
	}
}

// PublishSagaEvent отправляет событие SAGA (вызов из DialogService)
func (s *SagaOrchestrator) PublishSagaEvent(ctx context.Context, event domain.SagaEvent) error {
	log := logger.FromContext(ctx).With("transactionID", event.TransactionID, "eventType", event.Type)
	err := s.kafkaProducer.SendSagaEvent(ctx, event)
	if err != nil {
		log.Errorw("Failed to publish saga event", "err", err)
		return err
	}

	log.Infow("Saga event published")
	return nil
}

// HandleSagaEvent обрабатывает события SAGA и передает их в DialogService
func (s *SagaOrchestrator) HandleSagaEvent(ctx context.Context, event domain.SagaEvent) error {
	log := logger.FromContext(ctx).With("transactionID", event.TransactionID, "eventType", event.Type)

	switch event.Type {
	case domain.SagaCounterIncremented:
		log.Infow("Processing SagaCounterIncremented")
		return s.dialogService.CommitSagaTransaction(ctx, event.TransactionID)

	case domain.SagaCounterIncrementFailed:
		log.Infow("Processing SagaCounterIncrementFailed")
		return s.dialogService.RollbackSagaTransaction(ctx, event.TransactionID)
	case domain.SagaMessageSent:
		log.Debugw("Skipping SagaMessageSent")
		return nil
	default:
		log.Warnw("Unknown SAGA event type, ignoring")
		return nil
	}
}
