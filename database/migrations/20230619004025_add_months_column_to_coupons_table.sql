-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    coupons
ADD
    `months` INT DEFAULT 0
after
    `coupon`;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
`coupons` DROP `months`;

-- +goose StatementEnd