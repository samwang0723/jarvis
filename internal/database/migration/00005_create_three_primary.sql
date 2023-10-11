-- +goose Up
-- +goose StatementBegin
CREATE TABLE three_primary (
    id bigint unsigned NOT NULL PRIMARY KEY,
    stock_id varchar(8) NOT NULL,
    exchange_date varchar(32) NOT NULL,
    foreign_trade_shares bigint DEFAULT 0,
    trust_trade_shares bigint DEFAULT 0,
    dealer_trade_shares bigint DEFAULT 0,
    hedging_trade_shares bigint DEFAULT 0,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at datetime NULL,
    UNIQUE KEY index_stock_id_exchange_date(stock_id, exchange_date),
    KEY index_foreign(foreign_trade_shares),
    KEY index_trust(trust_trade_shares),
    KEY index_dealer(dealer_trade_shares),
    KEY index_hedging(hedging_trade_shares),
    KEY index_exchange_date(exchange_date),
    KEY index_stock_id(stock_id)
) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS three_primary;
-- +goose StatementEnd
