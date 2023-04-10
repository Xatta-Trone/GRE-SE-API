-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    lists
ADD CONSTRAINT `fk_lists_relation_folder_id` FOREIGN KEY(`folder_id`) REFERENCES `folders`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
-- +goose StatementEnd
