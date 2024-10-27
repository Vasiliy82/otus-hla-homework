-- +goose Up
-- +goose StatementBegin

CREATE TABLE posts (
    id BIGSERIAL NOT NULL,                             -- Уникальный идентификатор поста
    user_id UUID NOT NULL,                             -- ID пользователя-владельца поста
    message TEXT NOT NULL,                             -- Текстовое сообщение поста
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,  -- Время создания поста
    modified_at TIMESTAMPTZ,                           -- Время последнего редактирования (может быть NULL)
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users_friends

-- +goose StatementEnd
