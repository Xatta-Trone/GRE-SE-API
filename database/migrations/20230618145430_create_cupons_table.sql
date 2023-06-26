-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS coupons(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `coupon` VARCHAR(255) UNIQUE not null,
    `user_id` BIGINT UNSIGNED DEFAULT null,
    `expires` DATETIME DEFAULT null,
    `used` INT DEFAULT 0,
    `max_use` INT DEFAULT null,
    CONSTRAINT `fk_coupons_user_id` FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    PRIMARY KEY `pk_id`(`id`)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `coupons`;

-- +goose StatementEnd