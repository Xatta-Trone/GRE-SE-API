-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `word_groups` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255),
    `created_at` DATETIME NULL,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY `pk_id`(`id`)
) ENGINE = InnoDB;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE word_groups;

-- +goose StatementEnd