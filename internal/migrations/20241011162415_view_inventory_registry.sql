-- +goose Up
CREATE VIEW inventory_item_with_info AS
SELECT
    inventory.*,
    registry.name,
    registry.GTIN,
    registry.category,
    registry.description,
    registry.unit_type,
    registry.quantity,
    registry.expiration_after_opening,
    registry.nutrients_per_item,
    registry.calories,
    registry.fats,
    registry.saturated,
    registry.carbs,
    registry.sugars,
    registry.protein,
    registry.fiber,
    registry.sodium
FROM inventory inventory
JOIN food_registry registry ON inventory.product_id = registry.product_id;

-- +goose Down
DROP VIEW IF EXISTS inventory_item_with_info;