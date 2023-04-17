-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `saved_lists`(
    `user_id` BIGINT UNSIGNED NOT NULL,
    `list_id` BIGINT UNSIGNED NOT null,
    CONSTRAINT `fk_saved_lists_user_id` FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_saved_lists_list_id` FOREIGN KEY(`list_id`) REFERENCES `lists`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS saved_lists;

-- +goose StatementEnd