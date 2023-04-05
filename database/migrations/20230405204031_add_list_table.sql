-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `lists`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED,
    `list_meta_id` BIGINT UNSIGNED,
    `name` VARCHAR(255) NOT NULL DEFAULT "Unnamed",
    `slug` VARCHAR(255) NOT NULL UNIQUE,
    `visibility` TINYINT(1) NOT NULL DEFAULT "1",
    `status` TINYINT(1) NOT NULL DEFAULT "0",
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY `pk_id`(`id`),
    CONSTRAINT `fk_lists_relation_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_lists_relation_list_meta_id` FOREIGN KEY (`list_meta_id`) REFERENCES `list_meta` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE = INNODB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS `lists`;
-- +goose StatementEnd
