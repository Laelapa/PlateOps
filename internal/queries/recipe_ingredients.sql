-- name: CreateRecipeIngredient :one
INSERT INTO recipe_ingredients (
    recipe_id,
    product_id,
    quantity
) VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetRecipeIngredientsByRecipe :many
SELECT *
FROM recipe_ingredients
WHERE recipe_id = $1;

-- name: GetRecipeIngredientsByProduct :many
SELECT *
FROM recipe_ingredients
WHERE product_id = $1;

-- name: GetRecipeIngredient :one
SELECT *
FROM recipe_ingredients
WHERE id = $1;

-- name: UpdateRecipeIngredient :exec
UPDATE recipe_ingredients
SET quantity = $2
WHERE id = $1;

-- name: DeleteRecipeIngredient :exec
DELETE FROM recipe_ingredients
WHERE id = $1;
