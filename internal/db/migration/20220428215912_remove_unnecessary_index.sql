-- +goose Up
-- +goose StatementBegin
ALTER TABLE three_primary 
    DROP INDEX `index_stock_id`,
    DROP INDEX `index_dealer`,
    DROP INDEX `index_hedging`;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE three_primary 
    ADD INDEX index_stock_id(stock_id),
    ADD INDEX index_dealer(dealer_trade_shares),
    ADD INDEX index_hedging(hedging_trade_shares);
-- +goose StatementEnd
