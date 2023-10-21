-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  users (
    id bigint unsigned NOT NULL PRIMARY KEY,
    first_name varchar(64) NOT NULL,
    last_name varchar(64) NOT NULL,
    email varchar(255) NOT NULL,
    phone varchar(128) NOT NULL,
    password char(60) NOT NULL,
    session_id VARCHAR(128),
    email_confirmed_at datetime DEFAULT NULL,
    phone_confirmed_at datetime DEFAULT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    session_expired_at datetime DEFAULT NULL,
    deleted_at datetime DEFAULT NULL,
    UNIQUE KEY index_email (email),
    UNIQUE KEY index_phone (phone)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
-- +goose StatementEnd
