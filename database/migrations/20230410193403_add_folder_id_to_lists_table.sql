-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    lists
ADD
    `folder_id` BIGINT UNSIGNED Null after `list_meta_id`;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `lists` DROP `folder_id`;

-- +goose StatementEnd