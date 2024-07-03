DROP INDEX IF EXISTS index_category;

ALTER TABLE stocks
    DROP COLUMN IF EXISTS category;
