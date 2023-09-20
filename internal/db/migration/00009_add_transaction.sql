-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  transactions (
    id bigint unsigned NOT NULL PRIMARY KEY,
    user_id bigint unsigned NOT NULL,
    reference_id bigint unsigned DEFAULT NULL,
    stock_id varchar(8) NOT NULL,
    order_type int DEFAULT 0,
    trade_price DECIMAL(8, 2) NOT NULL,
    quantity bigint unsigned DEFAULT 0,
    credit_amount DECIMAL(8, 2) DEFAULT 0.0,
    debit_amount DECIMAL(8, 2) DEFAULT 0.0,
    exchange_date varchar(32) NOT NULL,
    description varchar(255),
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at datetime DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    KEY index_reference_id (reference_id),
    KEY index_order_type (order_type),
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;

-- +goose StatementEnd
