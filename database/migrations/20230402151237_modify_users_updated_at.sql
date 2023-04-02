-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    `users`
MODIFY
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
