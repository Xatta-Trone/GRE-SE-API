-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    `wordlist`
MODIFY
    `word` VARCHAR(255) NOT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
-- +goose StatementEnd