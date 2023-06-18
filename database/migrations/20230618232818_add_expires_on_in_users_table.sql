-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    users
ADD
    `expires_on` DATETIME DEFAULT NULL
after
    `username`;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `users` DROP `expires_on`;

-- +goose StatementEnd