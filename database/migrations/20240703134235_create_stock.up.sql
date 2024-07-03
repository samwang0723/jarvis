CREATE TABLE stocks (
    id BIGINT NOT NULL PRIMARY KEY,
    stock_id VARCHAR(8) NOT NULL,
    name VARCHAR(32) NOT NULL,
    country VARCHAR(2) NOT NULL,
    site VARCHAR(16),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE (stock_id, country)
);

-- Create a trigger function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER update_updated_at BEFORE UPDATE ON stocks
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

