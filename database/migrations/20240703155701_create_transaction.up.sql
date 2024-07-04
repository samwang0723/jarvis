BEGIN;

CREATE TABLE transactions (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL,
    order_id uuid NOT NULL,
    order_type varchar(32) NOT NULL,
    credit_amount money NOT NULL DEFAULT 0.0,
    debit_amount money NOT NULL DEFAULT 0.0,
    status varchar(32) NOT NULL,
    version integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create an index on order_id
CREATE INDEX index_transaction_order_id ON transactions(order_id);
CREATE INDEX index_transaction_user_id ON transactions(user_id);

-- Create a trigger to call the function before update
CREATE TRIGGER update_transactions_updated_at
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
