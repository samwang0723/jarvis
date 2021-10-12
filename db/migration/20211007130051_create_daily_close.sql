-- +goose Up
-- +goose StatementBegin
CREATE TABLE daily_closes (
    id bigint unsigned NOT NULL PRIMARY KEY,
    stock_id varchar(8) NOT NULL,
    exchange_date varchar(32) NOT NULL,
    trade_shares bigint unsigned DEFAULT 0,
    transactions bigint unsigned DEFAULT 0,
    turnover bigint unsigned DEFAULT 0,
    open DECIMAL(8,2) NOT NULL,
    close DECIMAL(8,2) NOT NULL,
    high DECIMAL(8,2) NOT NULL,
    low DECIMAL(8,2) NOT NULL,
    price_diff DECIMAL(6,2) NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at datetime NULL,
    UNIQUE KEY index_stock_id_exchange_date(stock_id, exchange_date),
    KEY index_transactions(trade_shares),
    KEY index_close(close)
) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS daily_closes;
-- +goose StatementEnd
