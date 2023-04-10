-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `folder_list_relation`(
    `folder_id` BIGINT UNSIGNED NOT NULL,
    `list_id` BIGINT UNSIGNED NOT null,
    CONSTRAINT `fk_folder_list_relation_word_id` FOREIGN KEY(`folder_id`) REFERENCES `folders`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_folder_list_relation_list_id` FOREIGN KEY(`list_id`) REFERENCES `lists`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS folder_list_relation;

-- +goose StatementEnd