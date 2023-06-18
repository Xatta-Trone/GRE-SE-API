-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS learning_status(
    `user_id` BIGINT UNSIGNED NOT NULL,
    `list_id` BIGINT UNSIGNED NOT null,
    `word_id` INT UNSIGNED NOT null,
    `learning_state` TINYINT(1) NOT NULL DEFAULT 0,
    CONSTRAINT `fk_learning_status_user_id` FOREIGN KEY(`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_learning_status_list_id` FOREIGN KEY(`list_id`) REFERENCES `lists`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_learning_status_word_id` FOREIGN KEY(`word_id`) REFERENCES `words`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    UNIQUE KEY `learning_status_user_list_word_unique_key` (`user_id`, `list_id`, `word_id`)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS learning_status;
-- +goose StatementEnd