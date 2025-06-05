-- +goose Up
CREATE TABLE food_registry (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    GTIN VARCHAR(14) UNIQUE,
    category VARCHAR(255),
    description TEXT,
    unit_type VARCHAR(50), -- defines the units used for the quantity metric
    quantity INT,
    portion_count INT DEFAULT 1,
    expiration_after_opening INT,
    nutrients_per_item BOOLEAN DEFAULT FALSE, -- if `false` then it's per 100g or 100ml, otherwise it's per 1 unit of whatever the unit type is
    calories REAL,
    fats REAL DEFAULT 0,
    saturated REAL DEFAULT 0,
    carbs REAL DEFAULT 0,
    sugars REAL DEFAULT 0,
    protein REAL DEFAULT 0,
    fiber REAL DEFAULT 0,
    sodium REAL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by uuid REFERENCES users(id),
    updated_by uuid REFERENCES users(id),
    CONSTRAINT unit_type_check CHECK (unit_type IN ('items', 'grams', 'ml', 'portions'))
);

-- +goose Down
DROP TABLE food_registry;