-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    id            UUID PRIMARY KEY,

    type          TEXT             NOT NULL,

    sell_currency TEXT             NOT NULL,
    sell_quantity DOUBLE PRECISION,
    sell_address  TEXT             NOT NULL,
    sell_user_id  TEXT             NOT NULL,

    price         DOUBLE PRECISION NOT NULL,

    buy_currency  TEXT             NOT NULL,
    buy_quantity  DOUBLE PRECISION,
    buy_address   TEXT             NOT NULL,
    buy_user_id   TEXT             NOT NULL,

    status        TEXT             NOT NULL,

    created_at    TIMESTAMPTZ        NOT NULL,
    updated_at    TIMESTAMPTZ        NOT NULL,
    removed_at    TIMESTAMPTZ,
    closed_at     TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
