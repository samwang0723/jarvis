-- +goose Up
-- +goose StatementBegin
ALTER TABLE stocks
    ADD COLUMN market varchar(10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stocks DROP COLUMN market;
-- +goose StatementEnd
