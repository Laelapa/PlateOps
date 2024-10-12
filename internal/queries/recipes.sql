-- name: CreateRecipe :one
INSERT INTO recipes (
    user_id, 
    name,
    private,
    description,
    instructions,
    unit_type,
    quantity
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetRecipesByUser :many
SELECT *
FROM recipes
WHERE user_id = $1;

-- name: GetRecipeById :one
SELECT *
FROM recipes
WHERE id = $1;

-- name: GetRecipesByUserAndName :many
SELECT *
FROM recipes
WHERE user_id = $1 AND name = $2;

-- name: UpdateRecipe :exec
UPDATE recipes
SET
    user_id = $2,
    name = $3,
    private = $4,
    description = $5,
    instructions = $6,
    unit_type = $7,
    quantity = $8,
    nutrients_per_portion = $9,
    calories = $10,
    fats = $11,
    saturated = $12,
    carbs = $13,
    sugars = $14,
    protein = $15,
    fiber = $16,
    sodium = $17,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- should also call DeleteRecipeIngredients for the specific recipe id
-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = $1;
