-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `words` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `word` VARCHAR(255),
    `word_data` JSON NULL,
    `created_at` DATETIME NULL,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY `pk_id`(`id`),
    CONSTRAINT words_word_unique UNIQUE (word)
) ENGINE = InnoDB;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `words` DROP INDEX `words_word_unique`;

DROP TABLE `words`;

-- +goose StatementEnd