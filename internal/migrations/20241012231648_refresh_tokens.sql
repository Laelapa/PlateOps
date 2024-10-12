-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY UNIQUE,
    user_id uuid REFERENCES users(id), -- no cascade deletion for security auditing reasons
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP DEFAULT NULL, -- invalidation by the server
    revocation_reason TEXT DEFAULT NULL
    logged_out_at TIMESTAMP DEFAULT NULL, -- invalidation by the user
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,
);

-- +goose Down
DROP TABLE refresh_tokens;