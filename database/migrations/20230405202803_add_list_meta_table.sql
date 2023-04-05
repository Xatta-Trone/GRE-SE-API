-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `list_meta`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED,
    `name` VARCHAR(255) NOT NULL DEFAULT "Unnamed",
    `url` VARCHAR(255) NULL,
    `words` MEDIUMTEXT NULL,
    `visibility` TINYINT(1) NOT NULL DEFAULT "1",
    `status` TINYINT(1) NOT NULL DEFAULT "0",
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY `pk_id`(`id`),
    CONSTRAINT `fk_list_meta_relation_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE = INNODB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS `list_meta`;
-- +goose StatementEnd
