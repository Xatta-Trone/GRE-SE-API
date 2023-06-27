-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    coupons
MODIFY
    COLUMN `max_use` INT DEFAULT 0;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd