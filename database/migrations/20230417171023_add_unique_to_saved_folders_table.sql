-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    saved_folders
ADD
    CONSTRAINT saved_folders_unique UNIQUE (user_id, folder_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    saved_folders DROP CONSTRAINT saved_folders_unique;

-- +goose StatementEnd