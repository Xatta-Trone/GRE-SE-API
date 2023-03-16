-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    `adsf`
ADD
    `asdfasdf` INT NOT NULL
AFTER
    `asdf`;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- SELECT
--     'down SQL query';

-- +goose StatementEnd