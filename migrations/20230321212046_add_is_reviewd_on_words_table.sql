-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    words
ADD
    is_reviewed TINYINT(1) NOT NULL DEFAULT 0
AFTER
    word;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    words DROP is_reviewed;

-- +goose StatementEnd