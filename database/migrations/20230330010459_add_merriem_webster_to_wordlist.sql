-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    `mw` JSON DEFAULT NULL,
ADD
    `is_parsed_mw` TINYINT NOT NULL DEFAULT 0,
ADD
    `mw_try` INT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `wordlist` DROP `mw`,
    DROP `is_parsed_mw`,
    DROP `mw_try`;

-- +goose StatementEnd