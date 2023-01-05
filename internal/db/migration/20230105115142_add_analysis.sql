-- +goose Up
-- +goose StatementBegin
ALTER TABLE daily_closes
    ADD COLUMN half_year_high DECIMAL(8,2),
    ADD COLUMN average_fivedays_volume bigint unsigned,
    ADD COLUMN above_all_ma BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE daily_closes
    DROP COLUMN half_year_high,
    DROP COLUMN average_fivedays_volume,
    DROP COLUMN above_all_ma;
-- +goose StatementEnd
