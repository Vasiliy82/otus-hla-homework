-- +goose Up

    ALTER TABLE messages ADD COLUMN transaction_id TEXT UNIQUE;
    ALTER TABLE messages ADD COLUMN saga_status TEXT DEFAULT 'pending';

-- +goose Down

    ALTER TABLE messages DROP COLUMN transaction_id;
    ALTER TABLE messages DROP COLUMN saga_status;
