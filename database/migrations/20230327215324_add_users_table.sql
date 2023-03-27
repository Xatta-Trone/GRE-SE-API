-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255),
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY `pk_id`(`id`)
) ENGINE = InnoDB;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;

-- +goose StatementEnd