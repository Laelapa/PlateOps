-- +goose Up
CREATE TABLE recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    product_id INT REFERENCES food_registry(product_id) ON DELETE CASCADE,
    quantity FLOAT NOT NULL
);

-- +goose Down
DROP TABLE recipe_ingredients;