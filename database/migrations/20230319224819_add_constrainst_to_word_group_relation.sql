-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    `word_group_relation`
ADD
    CONSTRAINT `fk_word_group_relation_word_id` FOREIGN KEY (`word_id`) REFERENCES `words` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
ADD
    CONSTRAINT `fk_word_group_relation_word_group_id` FOREIGN KEY (`word_group_id`) REFERENCES `word_groups` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    word_group_relation DROP CONSTRAINT fk_word_group_relation_word_id,
    DROP CONSTRAINT fk_word_group_relation_word_group_id;

-- +goose StatementEnd