-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    `word_groups`
ADD
    `file_name` VARCHAR(255) NULL
AFTER
    name,
ADD
    `words` MEDIUMTEXT NULL DEFAULT NULL
AFTER
    `file_name`,
ADD
    `new_words` MEDIUMTEXT NULL DEFAULT NULL
AFTER
    name,
ADD
    `status` TINYINT(1) NOT NULL DEFAULT 0
AFTER
    new_words;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `word_groups` DROP `file_name`,
    DROP `words`,
    drop `new_words`,
    DROP `status`;

-- +goose StatementEnd