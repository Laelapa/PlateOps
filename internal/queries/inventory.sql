-- name: CreateInventoryItem :one
INSERT INTO inventory (
    user_id,
    product_id,
    expiration_date,
    current_quantity
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- inventory_item_with_info is a view that combines the data in the inventory table with the generic info for that product from the food_registry table

-- name: GetInventoryOfUser :many
SELECT *
FROM inventory_item_with_info
WHERE user_id = $1;

-- name: GetInventoryItemById :one
SELECT *
FROM inventory_item_with_info
WHERE id = $1;

-- name: GetInventoryItemsByName :many
SELECT *
FROM inventory_item_with_info
WHERE name = $1;

-- name: GetInventoryItemsByGtin :many
SELECT *
FROM inventory_item_with_info
WHERE GTIN = $1;

-- name: UpdateInventoryItem :exec
UPDATE inventory
SET 
    product_id = $2,
    expiration_date = $3,
    purchase_date = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: InventoryItemMarkOpened :exec
UPDATE inventory
SET 
    is_opened = TRUE, 
    opened_date = CURRENT_TIMESTAMP, 
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: InventoryItemMarkUnopened :exec
UPDATE inventory
SET
    is_opened = FALSE,
    opened_date = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: InventoryItemAlterQuantity :exec
UPDATE inventory
SET
    current_quantity = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteInventoryItem :exec
DELETE FROM inventory
WHERE id = $1;
