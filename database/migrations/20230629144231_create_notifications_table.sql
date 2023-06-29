-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notifications(
    `id` CHAR(26) NOT NULL,
    `content` text not null,
    `user_id` BIGINT UNSIGNED DEFAULT NULL,
    `url` Varchar(255) DEFAULT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `read_at` DATETIME DEFAULT null,
    CONSTRAINT `fk_notifications_user_id` FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    PRIMARY KEY `pk_id`(`id`)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `notifications`;

-- +goose StatementEnd