-- +goose Up

    -- Создание таблицы для хранения счетчиков непрочитанных сообщений в диалогах
    CREATE TABLE IF NOT EXISTS dialog_counters (
        dialog_id TEXT PRIMARY KEY,
        unread_count INTEGER NOT NULL DEFAULT 0
    );

    -- Индекс для ускорения поиска по dialog_id (не обязателен, если есть PRIMARY KEY)
    CREATE INDEX IF NOT EXISTS idx_dialog_id ON dialog_counters(dialog_id);

-- +goose Down

    DROP TABLE dialog_counters