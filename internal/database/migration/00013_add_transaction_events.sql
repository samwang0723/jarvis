-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  transaction_events (
    aggregate_id bigint unsigned NOT NULL,
    parent_id bigint unsigned NOT NULL,
    event_type varchar(50) NOT NULL,
    payload blob NOT NULL,
    version int NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (aggregate_id, version),
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_events;

-- +goose StatementEnd
