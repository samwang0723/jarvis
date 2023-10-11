-- +goose Up
-- +goose StatementBegin
CREATE TABLE stake_concentration (
    id bigint unsigned NOT NULL PRIMARY KEY,
    stock_id varchar(8) NOT NULL,
    exchange_date varchar(32) NOT NULL,
    sum_buy_shares bigint unsigned DEFAULT 0,
    sum_sell_shares bigint unsigned DEFAULT 0,
    avg_buy_price DECIMAL(8,2) NOT NULL,
    avg_sell_price DECIMAL(8,2) NOT NULL,
    concentration_1 DECIMAL(3,2) DEFAULT 0.0,
    concentration_5 DECIMAL(3,2) DEFAULT 0.0,
    concentration_10 DECIMAL(3,2) DEFAULT 0.0,
    concentration_20 DECIMAL(3,2) DEFAULT 0.0,
    concentration_60 DECIMAL(3,2) DEFAULT 0.0,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at datetime NULL,
    UNIQUE KEY index_stock_id_exchange_date(stock_id, exchange_date),
    KEY index_exchange_date(exchange_date),
    KEY index_stock_id(stock_id)
) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stake_concentration;
-- +goose StatementEnd
