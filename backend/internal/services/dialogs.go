package services

import (
	"context"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
)

type dialogService struct {
	cfg        *config.DialogServiceConfig
	repository domain.DialogRepository
}

// NewDialogService создает новый экземпляр сервиса
// func NewDialogService(cfg *config.DialogServiceConfig, repo domain.DialogRepository) domain.DialogService {
func NewDialogService(cfg *config.DialogServiceConfig, repo domain.DialogRepository) domain.DialogService {
	return &dialogService{cfg: cfg, repository: repo}
}

// SendMessage отправляет сообщение пользователю
func (s *dialogService) SendMessage(ctx context.Context, myId, partnerId domain.UserKey, message string) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName(), "myId", myId, "partnerId", partnerId, "message", message)

	err := s.repository.SaveMessage(ctx, myId, partnerId, message)
	if err != nil {
		log.Errorw("s.repository.SaveMessage() returned error", "err", err)
		return fmt.Errorf("failed to save message %w", err)
	}

	log.Info("success")
	return nil
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
