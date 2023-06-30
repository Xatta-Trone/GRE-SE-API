-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    `in_words` tinyint(1) not null DEFAULT 0
after
    `word`,
ADD
    `tried` tinyint(1) not null DEFAULT 0
after
    `word`;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    wordlist DROP COLUMN `in_words`,
    `tried`;

-- +goose StatementEnd