-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `list_word_relation`(
    `word_id` INT UNSIGNED NOT NULL,
    `list_id` BIGINT UNSIGNED NOT null,
    CONSTRAINT `fk_list_word_relation_word_id` FOREIGN KEY(`word_id`) REFERENCES `words`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `fk_list_word_relation_list_id` FOREIGN KEY(`list_id`) REFERENCES `lists`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS list_word_relation
-- +goose StatementEnd
