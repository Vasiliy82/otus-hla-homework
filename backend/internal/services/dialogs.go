package services

import (
	"context"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
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
	err := s.repository.SaveMessage(ctx, myId, partnerId, message)
	if err != nil {
		return fmt.Errorf("failed to save message %w", err)
	}

	return nil
}

// GetDialog получает диалог между двумя пользователями
func (s *dialogService) GetDialog(ctx context.Context, myId, partnerId domain.UserKey, limit, offset int) ([]domain.DialogMessage, error) {

	messages, err := s.repository.GetDialog(ctx, myId, partnerId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages %w", err)
	}

	return messages, nil
}

// GetDialog получает диалог между двумя пользователями
func (s *dialogService) GetDialogs(ctx context.Context, myId domain.UserKey, limit, offset int) ([]domain.Dialog, error) {

	messages, err := s.repository.GetDialogs(ctx, myId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages %w", err)
	}

	return messages, nil
}
