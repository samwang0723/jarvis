CREATE OR REPLACE FUNCTION before_insert_unique_check()
RETURNS TRIGGER AS $$
DECLARE
    duplicate_count INT;
BEGIN
    SELECT COUNT(*)
    INTO duplicate_count
    FROM picked_stocks
    WHERE user_id = NEW.user_id
      AND stock_id = NEW.stock_id
      AND deleted_at IS NULL;

    IF duplicate_count > 0 THEN
        RAISE EXCEPTION 'Duplicate entry for user_id, stock_id with NULL deleted_at';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_unique_check
BEFORE INSERT ON picked_stocks
FOR EACH ROW
EXECUTE FUNCTION before_insert_unique_check();

