-- +goose Up
CREATE TABLE food_registry (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    GTIN VARCHAR(14) UNIQUE,
    category VARCHAR(255),
    description TEXT,
    unit_type VARCHAR(50) NOT NULL, -- defines the units used for the quantity metric
    quantity INT,
    portion_count INT DEFAULT 1,
    expiration_after_opening INT,
    nutrients_per_item BOOLEAN DEFAULT FALSE, -- if `false` then it's per 100g
    calories DECIMAL(10, 2),
    fats DECIMAL(10, 2),
    saturated DECIMAL(10, 2),
    carbs DECIMAL(10, 2),
    sugars DECIMAL(10, 2),
    protein DECIMAL(10, 2),
    fiber DECIMAL(10, 2),
    sodium DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by uuid REFERENCES users(id),
    uodated_by uuid REFERENCES users(id),
    CONSTRAINT unit_type_check CHECK (unit_type IN ('items', 'grams', 'ml', 'portions'))
);

-- +goose Down
DROP TABLE food_registry;