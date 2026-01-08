-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories
(
    id                 uuid PRIMARY KEY                     DEFAULT uuidv7(),
    name               text                        NOT NULL,
    owner_id           uuid                        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    parent_category_id uuid,
    type               text                        NOT NULL,
    created_at         timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
