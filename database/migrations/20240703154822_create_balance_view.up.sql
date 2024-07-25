BEGIN;

CREATE TABLE balance_views (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    balance money NOT NULL DEFAULT 0.0,
    available money NOT NULL DEFAULT 0.0,
    pending money NOT NULL DEFAULT 0.0,
    version integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_balance_views_updated_at
BEFORE UPDATE ON balance_views
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
