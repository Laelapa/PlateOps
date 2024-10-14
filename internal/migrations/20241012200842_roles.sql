-- +goose Up
CREATE TABLE roles (
    id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'user',
    CONSTRAINT role_check CHECK (role IN ('user', 'admin', 'limited'))
);

-- +goose Down
DROP TABLE roles;