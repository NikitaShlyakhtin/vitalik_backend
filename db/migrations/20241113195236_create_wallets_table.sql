-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets
(
    address    TEXT PRIMARY KEY,

    user_id    TEXT             NOT NULL,

    currency   TEXT             NOT NULL,
    balance    DOUBLE PRECISION NOT NULL,

    created_at TIMESTAMPTZ        NOT NULL,
    updated_at TIMESTAMPTZ       NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
-- +goose StatementEnd
