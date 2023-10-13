-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  transactions (
    id bigint unsigned NOT NULL PRIMARY KEY,
    user_id bigint unsigned NOT NULL,
    reference_id bigint unsigned DEFAULT NULL,
    stock_id varchar(8) NOT NULL,
    order_type VARCHAR(32) NOT NULL,
    trade_price DECIMAL(8, 2) NOT NULL,
    quantity bigint unsigned NOT NULL DEFAULT 0,
    credit_amount DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    debit_amount DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    exchange_date varchar(32) NOT NULL,
    description varchar(255),
    status VARCHAR(32) NOT NULL,
    version int NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY index_reference_id (reference_id)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;

-- +goose StatementEnd
