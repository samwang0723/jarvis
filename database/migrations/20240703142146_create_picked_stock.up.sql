BEGIN;

CREATE TABLE picked_stocks (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL,
    stock_id varchar(8) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp NULL,
    CONSTRAINT fk_stock_id FOREIGN KEY (stock_id) REFERENCES stocks (id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create a partial unique index to enforce the constraint
CREATE UNIQUE INDEX idx_unique_active_picked_stock_per_user
ON picked_stocks (user_id, stock_id)
WHERE deleted_at IS NULL;

CREATE TRIGGER update_picked_stocks_updated_at
BEFORE UPDATE ON picked_stocks
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
