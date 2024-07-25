CREATE TABLE order_events (
    aggregate_id uuid NOT NULL,
    parent_id uuid NOT NULL,
    event_type varchar(50) NOT NULL,
    payload jsonb NOT NULL,
    version integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (aggregate_id, version)
);
