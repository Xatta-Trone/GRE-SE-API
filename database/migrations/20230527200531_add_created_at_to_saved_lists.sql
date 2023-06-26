-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    saved_lists
ADD
    `created_at` DATETIME DEFAULT NULL
;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `saved_lists` DROP `created_at`;

-- +goose StatementEnd