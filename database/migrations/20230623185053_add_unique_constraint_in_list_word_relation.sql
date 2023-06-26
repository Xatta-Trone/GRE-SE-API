-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    list_word_relation
ADD
    CONSTRAINT list_word_relation_list_id_word_id_unique UNIQUE (list_id, word_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    list_word_relation DROP CONSTRAINT list_word_relation_list_id_word_id_unique;

-- +goose StatementEnd