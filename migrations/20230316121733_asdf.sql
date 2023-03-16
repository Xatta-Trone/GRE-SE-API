-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    words_api_try INT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
alter table
    `wordlist` DROP COLUMN `words_api_try`;

-- +goose StatementEnd