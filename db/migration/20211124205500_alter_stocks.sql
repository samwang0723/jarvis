-- +goose Up
-- +goose StatementBegin
ALTER TABLE stocks 
    ADD COLUMN category varchar(16),
    ADD INDEX index_category(category);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stocks DROP COLUMN category;
-- +goose StatementEnd
