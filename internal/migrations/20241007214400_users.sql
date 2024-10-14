-- +goose Up
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    last_login TIMESTAMP,
    active_account BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;