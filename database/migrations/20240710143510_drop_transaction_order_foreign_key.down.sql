ALTER TABLE transactions ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
