-- +goose Up
-- +goose StatementBegin
CREATE TABLE blacklisted (
    serial BIGINT NOT NULL PRIMARY KEY
    , expire_date TIMESTAMP NOT NULL);

CREATE INDEX blacklisted_expire_date_idx ON blacklisted (expire_date);

CREATE SEQUENCE jwt_token AS BIGINT
    INCREMENT BY 1
    MINVALUE 1
    START WITH 1
    CACHE 100;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE jwt_token;

DROP TABLE blacklisted;
-- +goose StatementEnd
