-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `word_group_relation` (
    `word_id` INT UNSIGNED NOT NULL,
    `word_group_id` INT UNSIGNED NOT NULL,
    `created_at` DATETIME NULL,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE word_group_relation;


-- +goose StatementEnd