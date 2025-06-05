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
    name = COALESCE($2, name),
    GTIN = COALESCE($3, GTIN),
    category = COALESCE($4, category),
    description = COALESCE($5, description),
    unit_type = COALESCE($6, unit_type),
    quantity = COALESCE($7, quantity),
    portion_count = COALESCE($8, portion_count),
    expiration_after_opening = COALESCE($9, expiration_after_opening),
    nutrients_per_item = COALESCE($10, nutrients_per_item),
    calories = COALESCE($11, calories),
    fats = COALESCE($12, fats),
    saturated = COALESCE($13, saturated),
    carbs = COALESCE($14, carbs),
    sugars = COALESCE($15, sugars),
    protein = COALESCE($16, protein),
    fiber = COALESCE($17, fiber),
    sodium = COALESCE($18, sodium),
    updated_at = CURRENT_TIMESTAMP,
    updated_by = $19
WHERE product_id = $1;

-- name: DeleteFoodEntry :exec
DELETE FROM food_registry
WHERE product_id = $1;
