-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    wordlist
ADD
    google_try INT NOT NULL DEFAULT 0
;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    wordlist DROP COLUMN google_try;

-- +goose StatementEnd