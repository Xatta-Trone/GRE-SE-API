-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    list_meta
ADD
    `folder_id` BIGINT UNSIGNED NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `list_meta` DROP `list_meta`;

-- +goose StatementEnd