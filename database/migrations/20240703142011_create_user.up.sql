BEGIN;

CREATE TABLE users (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name varchar(64) NOT NULL,
    last_name varchar(64) NOT NULL,
    email varchar(255) NOT NULL,
    phone varchar(128) NOT NULL,
    password char(60) NOT NULL,
    session_id varchar(128),
    email_confirmed_at timestamp DEFAULT NULL,
    phone_confirmed_at timestamp DEFAULT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_timestamp,
    updated_at timestamp NOT NULL DEFAULT CURRENT_timestamp,
    session_expired_at timestamp DEFAULT NULL,
    deleted_at timestamp DEFAULT NULL,
    UNIQUE (email),
    UNIQUE (phone)
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
