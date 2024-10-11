-- name: CreateFoodEntry :one
INSERT INTO food_registry (
    name, 
    GTIN, 
    category, 
    description, 
    unit_type, 
    quantity,
    portion_count,
    expiration_after_opening,
    nutrients_per_item,
    calories,
    fats,
    saturated,
    carbs,
    sugars,
    protein,
    fiber,
    sodium,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
) RETURNING *;

-- name: GetFoodEntryById :one
SELECT *
FROM food_registry
WHERE product_id = $1;

-- name: GetFoodEntryByGtin :one
SELECT *
FROM food_registry
WHERE GTIN = $1;

-- name: GetFoodEntriesByName :many
SELECT *
FROM food_registry
WHERE name = $1;

-- name: GetFoodEntriesByCategory :many
SELECT *
FROM food_registry
WHERE category = $1;

-- name: UpdateFoodEntry :exec
UPDATE food_registry
SET
    name = $2,
    GTIN = $3,
    category = $4, 
    description = $5, 
    unit_type = $6, 
    quantity = $7,
    portion_count = $8,
    expiration_after_opening = $9,
    nutrients_per_item = $10,
    calories = $11,
    fats = $12,
    saturated = $13,
    carbs = $14,
    sugars = $15,
    protein = $16,
    fiber = $17,
    sodium = $18,
    updated_at = CURRENT_TIMESTAMP,
    updated_by = $19
WHERE product_id = $1;

-- name: DeleteFoodEntry :exec
DELETE FROM food_registry
WHERE product_id = $1;
