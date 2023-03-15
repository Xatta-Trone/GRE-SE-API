-- +goose Up
-- +goose StatementBegin
ALTER TABLE `wordlist` ADD UNIQUE(`word`)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `wordlist` DROP UNIQUE(`word`)
-- +goose StatementEnd
