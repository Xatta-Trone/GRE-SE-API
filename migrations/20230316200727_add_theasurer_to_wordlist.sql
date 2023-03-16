-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    `thesaurus` JSON DEFAULT NULL,
ADD
    `is_parsed_th` TINYINT NOT NULL DEFAULT 0,
ADD
    `th_try` INT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `wordlist` DROP `thesaurus`,
    DROP `is_parsed_th`,
    DROP `th_try`;

-- +goose StatementEnd