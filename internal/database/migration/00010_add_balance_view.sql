-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  balance_views (
    id bigint unsigned NOT NULL,
    balance DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    available DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    pending DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    version int NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balance_views;

-- +goose StatementEnd
