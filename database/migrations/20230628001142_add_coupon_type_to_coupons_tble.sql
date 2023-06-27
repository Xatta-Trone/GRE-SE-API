-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    coupons
ADD
    `type` varchar(255) not null DEFAULT "one_time"
after
    `coupon`;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    `coupons` DROP `type`;

-- +goose StatementEnd