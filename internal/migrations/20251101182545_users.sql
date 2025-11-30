-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id         uuid PRIMARY KEY     DEFAULT uuidv7(),
    name       text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS external_identities
(
    id          uuid PRIMARY KEY     DEFAULT uuidv7(),
    user_id     uuid        NOT NULL,
    provider    text        NOT NULL,
    external_id text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
