-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  orders (
    id bigint unsigned NOT NULL PRIMARY KEY,
    user_id bigint unsigned NOT NULL,
    stock_id varchar(8) NOT NULL,
    buy_price DECIMAL(8, 2) NOT NULL,
    buy_quantity bigint unsigned NOT NULL DEFAULT 0,
    buy_exchange_date varchar(32) NOT NULL,
    sell_price DECIMAL(8, 2) NOT NULL,
    sell_quantity bigint unsigned NOT NULL DEFAULT 0,
    sell_exchange_date varchar(32) NOT NULL,
    profit_loss DECIMAL(8, 2) NOT NULL,
    profitable_price DECIMAL(8, 2) NOT NULL,
    status VARCHAR(32) NOT NULL,
    version int NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY index_user_id (user_id),
    KEY index_user_stock_id (user_id, stock_id)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;

-- +goose StatementEnd
