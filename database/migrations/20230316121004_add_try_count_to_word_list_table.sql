-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    `wordlist`
ADD
    `wiki_try` INT NOT NULL DEFAULT 0;


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
alter table
    `wordlist` DROP COLUMN `wiki_try`;

-- +goose StatementEnd