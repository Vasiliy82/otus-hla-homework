package repository

import (
	"context"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
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
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("started", "myId", myId, "partnerId", partnerId, "message", message)

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		log.Debugw("r.dbCluster.GetDBPool() returned error", err)
		return fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	dialogId := r.getDialogId(myId, partnerId)
	log.Debugw("r.getDialogId()", "dialogId", dialogId)

	// Начало транзакции
	log.Debugw("Tran begin")
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			err1 := tx.Rollback(ctx) // Откат в случае паники
			_ = err1
			log.Debugw("Tran rollback (panic)")
			panic(p) // Переподнимаем панику
		} else if err != nil {
			err1 := tx.Rollback(ctx) // Откат в случае ошибки
			_ = err1
			log.Debugw("Tran rollback")
		} else {
			err1 := tx.Commit(ctx) // Фиксация транзакции
			_ = err1
			log.Debugw("Tran commit")
		}
	}()

	// Вставляем данные в `dialogs`
	query := `
		INSERT INTO dialogs (dialog_id, user_id)
		VALUES ($1, $2), ($1, $3)
		ON CONFLICT DO NOTHING;
	`
	log.Debugw("tx.Exec()", "query", query, "$1", dialogId, "$2", myId, "$3", partnerId)
	_, err = tx.Exec(ctx, query, dialogId, myId, partnerId)
	if err != nil {
		log.Debug("tx.Exec() returned error", "err", err)
		return fmt.Errorf("failed to create dialog: %w", err)
	}
	log.Debug("tx.Exec() finished")

	// Вставляем данные в `messages`
	query = `
		INSERT INTO messages (dialog_id, author_id, message)
		VALUES ($1, $2, $3);
	`

	log.Debugw("tx.Exec()", "query", query, "$1", dialogId, "$2", myId, "$3", message)
	_, err = tx.Exec(ctx, query, dialogId, myId, message)
	if err != nil {
		log.Debug("tx.Exec() returned error", "err", err)
		return fmt.Errorf("failed to insert message: %w", err)
	}
	log.Debug("tx.Exec() finished")

	// Успешная транзакция
	log.Debug("finished")
	return nil
}

// GetDialog получает все сообщения между двумя пользователями
func (r *dialogRepository) GetDialog(ctx context.Context, myId, partnerId domain.UserKey, limit, offset int) ([]domain.DialogMessage, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("started", "myId", myId, "partnerId", partnerId, "limit", limit, "offset", offset)

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		log.Debugw("r.dbCluster.GetDBPool() returned error", err)
		return nil, fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	dialogId := r.getDialogId(myId, partnerId)
	log.Debugw("r.getDialogId()", "dialogId", dialogId)

	query := `
		SELECT dialog_id, message_id, author_id, datetime, message
		FROM messages
		WHERE dialog_id = $1
		    AND saga_status = $2
		ORDER BY datetime DESC
		LIMIT $3 OFFSET $4;
	`

	log.Debugw("db.Query()", "query", query, "$1", dialogId, "$2", domain.TxCommitted, "$3", limit, "$4", offset)
	rows, err := db.Query(ctx, query, dialogId, domain.TxCommitted, limit, offset)
	if err != nil {
		log.Debugw("db.Query() returned error", "err", err)
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()
	log.Debugw("db.Query() sccess")

	var messages []domain.DialogMessage
	for rows.Next() {
		var msg domain.DialogMessage
		if err := rows.Scan(&msg.DialogId, &msg.MessageId, &msg.AuthorId, &msg.Datetime, &msg.Message); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Debugw("rows.Err() returned error", "err", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	log.Debug("finished")
	return messages, nil
}

// GetDialogs получает список диалогов для данного пользователя
func (r *dialogRepository) GetDialogs(ctx context.Context, myId domain.UserKey, limit, offset int) ([]domain.Dialog, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		log.Debugw("r.dbCluster.GetDBPool() returned error", err)
		return nil, fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	query := `
		SELECT user_id, dialog_id
		FROM dialogs
		WHERE user_id = $1
		ORDER BY dialog_id
		LIMIT $2 OFFSET $3;
	`

	log.Debugw("db.Query()", "query", query, "$1", myId, "$2", limit, "$3", offset)
	rows, err := db.Query(ctx, query, myId, limit, offset)
	if err != nil {
		log.Debugw("db.Query() returned error", "err", err)
		return nil, fmt.Errorf("failed to query dialogs: %w", err)
	}
	defer rows.Close()
	log.Debugw("db.Query() sccess")

	var dialogs []domain.Dialog
	for rows.Next() {
		var dlg domain.Dialog
		if err := rows.Scan(&dlg.UserId, &dlg.DialogId); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		dialogs = append(dialogs, dlg)
	}

	if err := rows.Err(); err != nil {
		log.Debugw("rows.Err() returned error", "err", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	log.Debug("finished")
	return dialogs, nil
}

// UpdateSagaStatus обновляет статус SAGA транзакции
func (r *dialogRepository) UpdateSagaStatus(ctx context.Context, transactionID string, status string) error {

	log := logger.FromContext(ctx).With("func", logger.GetFuncName())

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		log.Debugw("r.dbCluster.GetDBPool() returned error", err)
		return fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	query := `
		UPDATE messages
		SET saga_status = $1
		WHERE transaction_id = $2;`

	log.Debugw("db.Exec()", "query", query, "$1", status, "$2", transactionID)

	_, err = db.Exec(ctx, query, status, transactionID)
	if err != nil {
		log.Debug("db.Exec() returned error", "err", err)
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	log.Debug("db.Exec() finished")

	// Успешная транзакция
	log.Debug("finished")
	return nil
}

func (r *dialogRepository) SaveMessageWithSaga(ctx context.Context, myId, partnerId domain.UserKey, message, transactionID string) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	log.Debugw("started", "myId", myId, "partnerId", partnerId, "message", message, "transactionID", transactionID)

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		log.Errorw("r.dbCluster.GetDBPool() returned error", "err", err)
		return fmt.Errorf("failed to get pool from cluster: %w", err)
	}

	dialogId := r.getDialogId(myId, partnerId)
	log.Debugw("r.getDialogId()", "dialogId", dialogId)

	// Начало транзакции
	log.Debugw("Tran begin")
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx) // Откат в случае паники
			log.Debugw("Tran rollback (panic)")
			panic(p) // Переподнимаем панику
		} else if err != nil {
			_ = tx.Rollback(ctx) // Откат в случае ошибки
			log.Debugw("Tran rollback")
		} else {
			err = tx.Commit(ctx) // Фиксация транзакции
			if err != nil {
				log.Errorw("Tran commit failed", "err", err)
			} else {
				log.Debugw("Tran commit successful")
			}
		}
	}()

	// Вставляем данные в `dialogs`
	query := `
		INSERT INTO dialogs (dialog_id, user_id)
		VALUES ($1, $2), ($1, $3)
		ON CONFLICT DO NOTHING;
	`
	log.Debugw("tx.Exec()", "query", query, "$1", dialogId, "$2", myId, "$3", partnerId)
	_, err = tx.Exec(ctx, query, dialogId, myId, partnerId)
	if err != nil {
		log.Errorw("tx.Exec() failed to insert into dialogs", "err", err)
		return fmt.Errorf("failed to create dialog: %w", err)
	}

	// Вставляем данные в `messages` с `transaction_id` и `saga_status`
	query = `
		INSERT INTO messages (dialog_id, author_id, message, transaction_id, saga_status)
		VALUES ($1, $2, $3, $4, 'pending');
	`
	log.Debugw("tx.Exec()", "query", query, "$1", dialogId, "$2", myId, "$3", message, "$4", transactionID)
	_, err = tx.Exec(ctx, query, dialogId, myId, message, transactionID)
	if err != nil {
		log.Errorw("tx.Exec() failed to insert message", "err", err)
		return fmt.Errorf("failed to insert message: %w", err)
	}

	// Успешная транзакция
	log.Debug("finished")
	return nil
}
