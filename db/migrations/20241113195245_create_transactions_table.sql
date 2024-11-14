-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions
(
    id               UUID PRIMARY KEY,

    sender_address   TEXT             NOT NULL,
    sender_user_id   TEXT             NOT NULL,

    receiver_address TEXT             NOT NULL,
    receiver_user_id TEXT             NOT NULL,

    amount           DOUBLE PRECISION NOT NULL,
    currency         TEXT             NOT NULL,

    purpose          TEXT,

    created_at       TIMESTAMPTZ        NOT NULL,
    updated_at       TIMESTAMPTZ        NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
