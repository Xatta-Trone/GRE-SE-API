-- +goose Up
-- +goose StatementBegin
UPDATE
    `wordlist`
SET
    `tried` = '1',
    `in_words` = '1'
WHERE
    `word` IN (
        SELECT
            word
        from
            words
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
UPDATE
    `wordlist`
SET
    `tried` = '0',
    `in_words` = '0'
WHERE
    `word` IN (
        SELECT
            word
        from
            words
    );

-- +goose StatementEnd