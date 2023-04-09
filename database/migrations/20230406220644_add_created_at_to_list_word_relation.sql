-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    list_word_relation
ADD
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `list_word_relation` DROP `created_at`;

-- +goose StatementEnd