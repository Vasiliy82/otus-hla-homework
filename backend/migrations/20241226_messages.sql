-- +goose Up

-- Создаем таблицу dialogs
CREATE TABLE dialogs (
    user_id UUID NOT NULL,         -- Идентификатор автора сообщения
    dialog_id CHAR(73) NOT NULL,    -- идентификатор диалога
    PRIMARY KEY (user_id, dialog_id)
);

-- Создаем таблицу messages
CREATE TABLE messages (
    dialog_id CHAR(73) NOT NULL,
    message_id SERIAL NOT NULL,
    author_id UUID NOT NULL,         -- Идентификатор автора сообщения
    datetime TIMESTAMP NOT NULL DEFAULT NOW(),    -- Время сообщения
    message TEXT NOT NULL,          -- Текст сообщения
    PRIMARY KEY (dialog_id, message_id)
);

-- Настраиваем шардирование
-- SELECT create_distributed_table('dialogs', 'user_id');
-- SELECT create_distributed_table('messages', 'dialog_id');

-- Индексация для ускорения запросов
-- Индекс для запросов по user1_id с сортировкой по datetime
CREATE INDEX idx_user1_datetime ON messages (dialog_id, datetime DESC);

-- +goose Down

DROP TABLE messages
DROP TABLE dialogs

