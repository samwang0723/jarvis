-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  picked_stocks (
    id bigint unsigned NOT NULL PRIMARY KEY,
    user_id bigint unsigned NOT NULL,
    stock_id varchar(8) NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at datetime NULL,
    UNIQUE KEY index_user_id (user_id)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS picked_stocks;

-- +goose StatementEnd
