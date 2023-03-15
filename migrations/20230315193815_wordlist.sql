-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `wordlist` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `word` VARCHAR(255),
    `is_all_parsed` TINYINT(1) NOT NULL DEFAULT "0",
    `google` JSON DEFAULT NULL,
    `is_google_parsed` TINYINT(1) NOT NULL DEFAULT "0",
    `wiki` JSON DEFAULT NULL,
    `is_wiki_parsed` TINYINT(1) NOT NULL DEFAULT "0",
    `words_api` JSON DEFAULT NULL,
    `is_words_api_parsed` TINYINT(1) NOT NULL DEFAULT "0",
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY `pk_id`(`id`)
) ENGINE = InnoDB;

-- ALTER TABLE
--     `wordlist`
-- ADD
--     UNIQUE KEY `words_unique_key_name` (`word`);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `wordlist`;

-- +goose StatementEnd