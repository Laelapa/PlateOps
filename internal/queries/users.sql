-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserIdByEmail :one
SELECT id
FROM users
WHERE email = $1;

-- name: GetUserIdByUsername :one
SELECT id
FROM users
WHERE username = $1;

-- name: UpdateUser :exec
UPDATE users
SET username = $2, email = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserUsername :exec
UPDATE users
SET username = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: EnableUser :exec
UPDATE users
SET active_account=TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DisableUser :exec
UPDATE users
SET active_account=FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;