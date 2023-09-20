-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  balance_views (
    id bigint unsigned NOT NULL PRIMARY KEY,
    user_id bigint unsigned NOT NULL,
    current_amount DECIMAL(8, 2) DEFAULT 0.0,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at datetime DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balance_views;

-- +goose StatementEnd
