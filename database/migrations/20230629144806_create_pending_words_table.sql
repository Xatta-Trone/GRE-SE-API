-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pending_words(
    `word` VARCHAR(255) not null,
    `list_id` BIGINT UNSIGNED DEFAULT NULL,
    `approved` tinyint not null DEFAULT 0,
    CONSTRAINT `fk_pending_words_list_id` FOREIGN KEY(`list_id`) REFERENCES `lists`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pending_words`;

-- +goose StatementEnd