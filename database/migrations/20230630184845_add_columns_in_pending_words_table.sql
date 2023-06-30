-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    pending_words
ADD
    `tried` tinyint(1) not null DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    pending_words DROP COLUMN `tried`;

-- +goose StatementEnd