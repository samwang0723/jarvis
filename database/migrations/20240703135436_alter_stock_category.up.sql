ALTER TABLE stocks
    ADD COLUMN category VARCHAR(16);

CREATE INDEX index_category ON stocks (category);

