-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions
(
    id               uuid PRIMARY KEY                     DEFAULT uuidv7(),
    user_id          uuid                        NOT NULL,
    category_id      uuid                        NOT NULL,
    amount           numeric(10, 2)              NOT NULL,
    created_at       timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
