package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/tarantool/go-tarantool"
)

type dialogRepositoryTar struct {
	conn *tarantool.Connection
}

func NewdialogRepositoryTar(conn *tarantool.Connection) domain.DialogRepository {
	return &dialogRepositoryTar{conn: conn}
}

func (r *dialogRepositoryTar) SaveMessage(ctx context.Context, myId, partnerId domain.UserKey, message string) error {
	// Вызов Lua-функции save_message
	_, err := r.conn.Call("save_message", []interface{}{myId, partnerId, message})
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

func (r *dialogRepositoryTar) GetMessages(ctx context.Context, myId, partnerId domain.UserKey, limit, offset int) ([]domain.DialogMessage, error) {
	// Вызов Lua-функции get_messages
	resp, err := r.conn.Call("get_messages", []interface{}{myId, partnerId, limit, offset})
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Преобразуем результат
	var messages []domain.DialogMessage
	for _, tuple := range resp.Data {
		resp_arr := tuple.([]interface{})
		if len(resp_arr) == 0 {
			return nil, nil
		}
		datetimeFloat := resp_arr[4].(float64)
		datetime := time.Unix(int64(datetimeFloat), int64((datetimeFloat-float64(int64(datetimeFloat)))*1e9))

		msg := domain.DialogMessage{
			DialogId:  domain.DialogKey(resp_arr[0].(string)),
			MessageId: domain.MessageKey(resp_arr[1].(uint64)),
			AuthorId:  domain.UserKey(resp_arr[2].(string)),
			Message:   resp_arr[3].(string),
			Datetime:  datetime,
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *dialogRepositoryTar) GetDialogs(ctx context.Context, myId domain.UserKey, limit, offset int) ([]domain.Dialog, error) {
	// Вызов Lua-функции get_dialogs
	resp, err := r.conn.Call("get_dialogs", []interface{}{myId, limit, offset})
	if err != nil {
		return nil, fmt.Errorf("failed to get dialogs: %w", err)
	}

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

	return dialogs, nil
}
