-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token,
    user_id,
    expires_at,
    user_agent,
    ip_address
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTokensByUser :many
SELECT *
FROM refresh_tokens
WHERE user_id = $1;

-- name: GetTokensByIp :many
SELECT *
FROM refresh_tokens
WHERE ip_address = $1;

-- name: GetToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, revocation_reason = $2
WHERE token = $1;

-- name: RevokeTokensOfUser :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, revocation_reason = $2
WHERE user_id = $1;

-- name: RevokeTokensOfIp :exec
UPDATE refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP, revocation_reason = $2
WHERE ip_address = $1;

-- name: LogOutToken :exec
UPDATE refresh_tokens
SET logged_out_at = CURRENT_TIMESTAMP
WHERE token = $1;