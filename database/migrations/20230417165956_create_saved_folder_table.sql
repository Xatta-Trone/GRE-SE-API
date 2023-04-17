-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `saved_folders`(
    `user_id` BIGINT UNSIGNED NOT NULL,
    `folder_id` BIGINT UNSIGNED NOT null,
    CONSTRAINT `fk_saved_folders_user_id` FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_saved_folders_folder_id` FOREIGN KEY(`folder_id`) REFERENCES `folders`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS saved_folders;

-- +goose StatementEnd