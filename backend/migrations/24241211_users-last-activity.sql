-- +goose Up
-- +goose StatementBegin

CREATE TABLE users_last_activity (
    id UUID NOT NULL,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id) REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE INDEX users_last_activity_last_activity_idx ON users_last_activity(last_activity DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users_last_activity

-- +goose StatementEnd
