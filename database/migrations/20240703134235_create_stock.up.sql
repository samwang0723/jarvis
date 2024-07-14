BEGIN;

CREATE TABLE stocks (
    id varchar(8) NOT NULL PRIMARY KEY,
    name varchar(32) NOT NULL,
    country varchar(2) NOT NULL,
    category varchar(16) NOT NULL,
    market varchar(10) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp NULL,
    UNIQUE (id, country)
);

-- Create the trigger
CREATE TRIGGER update_stocks_updated_at BEFORE UPDATE ON stocks
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;

