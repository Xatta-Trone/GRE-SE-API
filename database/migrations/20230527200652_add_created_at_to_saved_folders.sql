-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    saved_folders
ADD
    `created_at` DATETIME DEFAULT NULL
;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `saved_folders` DROP `created_at`;

-- +goose StatementEnd