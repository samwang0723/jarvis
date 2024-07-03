CREATE TABLE balance_views (
    id BIGINT NOT NULL PRIMARY KEY,
    balance DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    available DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    pending DECIMAL(8, 2) NOT NULL DEFAULT 0.0,
    version INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create a trigger to update the updated_at column on row update
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_balance_views_updated_at
BEFORE UPDATE ON balance_views
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
