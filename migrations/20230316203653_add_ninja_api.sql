-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    `ninja` JSON DEFAULT NULL,
ADD
    `is_parsed_ninja` TINYINT NOT NULL DEFAULT 0,
ADD
    `ninja_try` INT NOT NULL DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `wordlist` DROP `ninja`,
    DROP `is_parsed_ninja`,
    DROP `ninja_try`;

-- +goose StatementEnd