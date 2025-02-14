package services

import (
	"context"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/utils"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
)

type dialogService struct {
	cfg        *config.DialogServiceConfig
	repository domain.DialogRepository
	sagaCoord  domain.SagaCoordinator
}

// NewDialogService создает новый экземпляр сервиса
// func NewDialogService(cfg *config.DialogServiceConfig, repo domain.DialogRepository) domain.DialogService {
func NewDialogService(cfg *config.DialogServiceConfig, repo domain.DialogRepository) domain.DialogService {
	return &dialogService{cfg: cfg, repository: repo}
}

// SendMessage отправляет сообщение пользователю
func (s *dialogService) SendMessage(ctx context.Context, myId, partnerId domain.UserKey, message string) (string, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName(), "myId", myId, "partnerId", partnerId, "message", message)

	transactionID := utils.GenerateID()

	err := s.repository.SaveMessageWithSaga(ctx, myId, partnerId, message, transactionID)
	if err != nil {
		log.Errorw("s.repository.SaveMessage() returned error", "err", err)
		return "", fmt.Errorf("failed to save message %w", err)
	}

	event := domain.SagaEvent{
		TransactionID: transactionID,
		Type:          domain.SagaMessageSent,
		DialogID:      s.getCountersDialogId(myId, partnerId),
	}

	if err = s.sagaCoord.PublishSagaEvent(ctx, event); err != nil {
		log.Warnw("Failure", "err", err)
		return "", err
	}

	log.Info("success")
	return transactionID, nil
}

// GetDialog получает диалог между двумя пользователями
func (s *dialogService) GetDialog(ctx context.Context, myId, partnerId domain.UserKey, limit, offset int) ([]domain.DialogMessage, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName(), "myId", myId, "partnerId", partnerId, "limit", limit, "offset", offset)

	dialog, err := s.repository.GetDialog(ctx, myId, partnerId, limit, offset)
	if err != nil {
		log.Errorw("s.repository.GetDialog() returned error", "err", err)
		return nil, fmt.Errorf("failed to retrieve dialog %w", err)
	}

	log.Info("success")
	return dialog, nil
}

// GetDialog получает диалог между двумя пользователями
func (s *dialogService) GetDialogs(ctx context.Context, myId domain.UserKey, limit, offset int) ([]domain.Dialog, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName(), "myId", myId, "limit", limit, "offset", offset)

	dialogs, err := s.repository.GetDialogs(ctx, myId, limit, offset)
	if err != nil {
		log.Errorw("s.repository.GetDialogs() returned error", "err", err)
		return nil, fmt.Errorf("failed to retrieve dialogs %w", err)
	}

	log.Info("success")
	return dialogs, nil
}

// Теперь `dialogService` реализует `SagaCoordinator`
func (s *dialogService) CommitSagaTransaction(ctx context.Context, transactionID string) error {
	return s.repository.UpdateSagaStatus(ctx, transactionID, domain.TxCommitted)
}

func (s *dialogService) RollbackSagaTransaction(ctx context.Context, transactionID string) error {
	return s.repository.UpdateSagaStatus(ctx, transactionID, domain.TxFailed)
}

func (s *dialogService) PublishSagaEvent(ctx context.Context, event domain.SagaEvent) error {
	return s.sagaCoord.PublishSagaEvent(ctx, event)
}

func (s *dialogService) SetSagaCoordinator(sagaCoordinator domain.SagaCoordinator) {
	s.sagaCoord = sagaCoordinator
}

func (s *dialogService) getCountersDialogId(myId, FriendId domain.UserKey) string {
	return fmt.Sprintf("%s:%s", FriendId, myId)
}
