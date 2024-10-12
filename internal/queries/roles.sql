-- name: CreateRole :one
INSERT INTO roles (
    id, role
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetRole :one
SELECT *
FROM roles
WHERE id = $1;

-- name: UpdateRole :exec
UPDATE roles
SET role = $2
WHERE id = $1;

-- Cannot delete role without deleting user