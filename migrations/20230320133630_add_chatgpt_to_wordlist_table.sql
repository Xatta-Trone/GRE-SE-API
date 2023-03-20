-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    `gpt` JSON DEFAULT NULL,
ADD
    `is_parsed_gpt` TINYINT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `wordlist` DROP `gpt`,
    DROP `is_parsed_gpt`;

-- +goose StatementEnd