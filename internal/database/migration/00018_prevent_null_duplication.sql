DELIMITER //

CREATE TRIGGER before_insert_unique_check BEFORE INSERT ON picked_stocks
FOR EACH ROW
BEGIN
  DECLARE duplicate_count INT;

  SELECT COUNT(*)
  INTO duplicate_count
  FROM picked_stocks
  WHERE user_id = NEW.user_id
    AND stock_id = NEW.stock_id
    AND (deleted_at IS NULL);

  IF duplicate_count > 0 THEN
    SIGNAL SQLSTATE '45000'
    SET MESSAGE_TEXT = 'Duplicate entry for user_id, stock_id with NULL deleted_at';
  END IF;
END;

//
DELIMITER ;

