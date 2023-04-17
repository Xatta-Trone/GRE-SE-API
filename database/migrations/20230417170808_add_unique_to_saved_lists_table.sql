-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    saved_lists
ADD
    CONSTRAINT saved_lists_unique UNIQUE (user_id, list_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    saved_lists DROP CONSTRAINT saved_lists_unique;

-- +goose StatementEnd