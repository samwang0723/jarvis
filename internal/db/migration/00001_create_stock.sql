-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  stocks (
    id bigint unsigned NOT NULL PRIMARY KEY,
    stock_id varchar(8) NOT NULL,
    name varchar(32) NOT NULL,
    country varchar(2) NOT NULL,
    site varchar(16),
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at datetime NULL,
    UNIQUE KEY index_stock_id (stock_id, country)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks;

-- +goose StatementEnd
