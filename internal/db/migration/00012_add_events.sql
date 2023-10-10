-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  events (
    id bigint unsigned NOT NULL PRIMARY KEY,
    aggregate_id bigint unsigned NOT NULL,
    event_type varchar(255) NOT NULL,
    payload json NOT NULL,
    version int NOT NULL,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (aggregate_id) REFERENCES transactions (id),
    KEY index_aggregate_id (aggregate_id)
  ) DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;

-- +goose StatementEnd
