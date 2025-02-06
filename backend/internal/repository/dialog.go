package repository

import (
	"context"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
)

type dialogRepository struct {
	dbCluster *postgresqldb.DBCluster
}

func NewDialogRepository(dbCluster *postgresqldb.DBCluster) domain.DialogRepository {
	return &dialogRepository{dbCluster: dbCluster}
}

func (r *dialogRepository) getDialogId(myId, partnerId domain.UserKey) string {
	if myId > partnerId {
		return fmt.Sprintf("%s:%s", partnerId, myId)
	}
	return fmt.Sprintf("%s:%s", myId, partnerId)
}

// SaveMessage сохраняет сообщение в таблицу
func (r *dialogRepository) SaveMessage(ctx context.Context, myId, partnerId domain.UserKey, message string) error {
	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	dialogId := r.getDialogId(myId, partnerId)

	// Начало транзакции
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			err1 := tx.Rollback(ctx) // Откат в случае паники
			_ = err1
			panic(p) // Переподнимаем панику
		} else if err != nil {
			err1 := tx.Rollback(ctx) // Откат в случае ошибки
			_ = err1

		} else {
			err1 := tx.Commit(ctx) // Фиксация транзакции
			_ = err1
		}
	}()

	// Вставляем данные в `dialogs`
	query := `
		INSERT INTO dialogs (dialog_id, user_id)
		VALUES ($1, $2), ($1, $3)
		ON CONFLICT DO NOTHING;
	`
	_, err = tx.Exec(ctx, query, dialogId, myId, partnerId)
	if err != nil {
		return fmt.Errorf("failed to create dialog: %w", err)
	}

	// Вставляем данные в `messages`
	query = `
		INSERT INTO messages (dialog_id, author_id, message)
		VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(ctx, query, dialogId, myId, message)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	// Успешная транзакция
	return nil
}

// GetMessages получает все сообщения между двумя пользователями
func (r *dialogRepository) GetMessages(ctx context.Context, myId, partnerId domain.UserKey, limit, offset int) ([]domain.DialogMessage, error) {

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	dialogId := r.getDialogId(myId, partnerId)

	query := `
		SELECT dialog_id, message_id, author_id, datetime, message
		FROM messages
		WHERE dialog_id = $1
		ORDER BY datetime DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := db.Query(ctx, query, dialogId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.DialogMessage
	for rows.Next() {
		var msg domain.DialogMessage
		if err := rows.Scan(&msg.DialogId, &msg.MessageId, &msg.AuthorId, &msg.Datetime, &msg.Message); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

// GetDialogs получает список диалогов для данного пользователя
func (r *dialogRepository) GetDialogs(ctx context.Context, myId domain.UserKey, limit, offset int) ([]domain.Dialog, error) {

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	query := `
		SELECT user_id, dialog_id
		FROM dialogs
		WHERE user_id = $1
		ORDER BY dialog_id
		LIMIT $2 OFFSET $3;
	`

	rows, err := db.Query(ctx, query, myId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query dialogs: %w", err)
	}
	defer rows.Close()

	var dialogs []domain.Dialog
	for rows.Next() {
		var dlg domain.Dialog
		if err := rows.Scan(&dlg.UserId, &dlg.DialogId); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		dialogs = append(dialogs, dlg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return dialogs, nil
}
