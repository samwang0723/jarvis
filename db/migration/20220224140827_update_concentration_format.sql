-- +goose Up
-- +goose StatementBegin
ALTER TABLE stake_concentration 
    MODIFY concentration_1 DECIMAL(5,2),
    MODIFY concentration_5 DECIMAL(5,2),
    MODIFY concentration_10 DECIMAL(5,2),
    MODIFY concentration_20 DECIMAL(5,2),
    MODIFY concentration_60 DECIMAL(5,2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stake_concentration 
    MODIFY concentration_1 DECIMAL(3,2),
    MODIFY concentration_5 DECIMAL(3,2),
    MODIFY concentration_10 DECIMAL(3,2),
    MODIFY concentration_20 DECIMAL(3,2),
    MODIFY concentration_60 DECIMAL(3,2);
-- +goose StatementEnd
