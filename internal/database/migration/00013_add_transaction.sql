-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  transactions (
    id bigint unsigned NOT NULL PRIMARY KEY,
    user_id bigint unsigned NOT NULL,
    order_id bigint unsigned NOT NULL,
    order_type VARCHAR(32) NOT NULL,
    credit_amount DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    debit_amount DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    status VARCHAR(32) NOT NULL,
    version int NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY index_order_id (order_id)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;

-- +goose StatementEnd
