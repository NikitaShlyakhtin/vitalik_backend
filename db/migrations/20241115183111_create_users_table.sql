-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    user_id TEXT PRIMARY KEY,
    hashed_password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
