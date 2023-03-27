-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    needs_attention TINYINT(1) NOT NULL DEFAULT 0
AFTER
    is_all_parsed;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    wordlist DROP needs_attention;

-- +goose StatementEnd