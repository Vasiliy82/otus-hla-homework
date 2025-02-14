package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/common/infrastructure/observability/logger"
	"github.com/tarantool/go-tarantool"
)

type dialogRepositoryTar struct {
	conn *tarantool.Connection
}

func NewdialogRepositoryTar(conn *tarantool.Connection) domain.DialogRepository {
	return &dialogRepositoryTar{conn: conn}
}

func (r *dialogRepositoryTar) SaveMessage(ctx context.Context, myId, partnerId domain.UserKey, message string) error {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	// Вызов Lua-функции save_message
	log.Debug("Calling tarantool: save_message")
	_, err := r.conn.Call("save_message", []interface{}{myId, partnerId, message})
	if err != nil {
		log.Debug("Calling tarantool: save_message error", "err", err)
		return fmt.Errorf("failed to save message: %w", err)
	}
	log.Debug("Calling tarantool: save_message success")
	return nil
}

func (r *dialogRepositoryTar) GetDialog(ctx context.Context, myId, partnerId domain.UserKey, limit, offset int) ([]domain.DialogMessage, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	// Вызов Lua-функции get_dialog
	log.Debug("Calling tarantool: get_dialog")
	resp, err := r.conn.Call("get_dialog", []interface{}{myId, partnerId, limit, offset})
	if err != nil {
		log.Debug("Calling tarantool: get_dialog error", "err", err)
		return nil, fmt.Errorf("failed to get dialog: %w", err)
	}
	log.Debug("Calling tarantool: get_dialog success")

	// Преобразуем результат
	var dialog []domain.DialogMessage
	for _, tuple := range resp.Data {
		resp_arr := tuple.([]interface{})
		if len(resp_arr) == 0 {
			return nil, nil
		}

		var datetime time.Time

		switch v := resp_arr[4].(type) {
		case float64:
			// Если значение float64, преобразуем к int64
			dtInt64 := int64(v)
			datetime = time.Unix(dtInt64, int64((v-float64(dtInt64))*1e9))
		case uint64:
			// Если значение uint64, преобразуем напрямую
			dtInt64 := int64(v)
			datetime = time.Unix(dtInt64, 0)

		default:
			log.Warnw("Unexpected type for datetime", "type", fmt.Sprintf("%T", v))
			datetime = time.Unix(0, 0)
		}

		msg := domain.DialogMessage{
			DialogId:  domain.DialogKey(resp_arr[0].(string)),
			MessageId: domain.MessageKey(resp_arr[1].(uint64)),
			AuthorId:  domain.UserKey(resp_arr[2].(string)),
			Message:   resp_arr[3].(string),
			Datetime:  datetime,
		}
		dialog = append(dialog, msg)
	}
	log.Debug("Tarantool get_dialog response successfully parsed")
	return dialog, nil
}

func (r *dialogRepositoryTar) GetDialogs(ctx context.Context, myId domain.UserKey, limit, offset int) ([]domain.Dialog, error) {
	log := logger.FromContext(ctx).With("func", logger.GetFuncName())
	// Вызов Lua-функции get_dialogs
	log.Debug("Calling tarantool: get_dialogs")
	resp, err := r.conn.Call("get_dialogs", []interface{}{myId, limit, offset})
	if err != nil {
		log.Debug("Calling tarantool: get_dialogs error", "err", err)
		return nil, fmt.Errorf("failed to get dialogs: %w", err)
	}
	log.Debug("Calling tarantool: get_dialogs success")

	// Преобразуем результат
	var dialogs []domain.Dialog
	for _, tuple := range resp.Data {
		resp_arr := tuple.([]interface{})
		if len(resp_arr) == 0 {
			return nil, nil
		}
		dlg := domain.Dialog{
			UserId:   domain.UserKey(resp_arr[0].(string)),
			DialogId: domain.DialogKey(resp_arr[1].(string)),
		}
		dialogs = append(dialogs, dlg)
	}
	log.Debug("Tarantool get_dialogs response successfully parsed")
	return dialogs, nil
}

func (r *dialogRepositoryTar) SaveMessageWithSaga(ctx context.Context, myId, partnerId domain.UserKey, message, transactionID string) error {
	return fmt.Errorf("Not implemented")
}

func (r *dialogRepositoryTar) UpdateSagaStatus(ctx context.Context, transactionID string, status string) error {
	return fmt.Errorf("Not implemented")
}
