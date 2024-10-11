-- +goose Up
CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    user_id uuid REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    private BOOLEAN DEFAULT TRUE,
    description TEXT,
    instructions TEXT,
    unit_type VARCHAR(50) NOT NULL, -- defines the units used for the quantity metric
    quantity INT,
    nutrients_per_portion BOOLEAN DEFAULT FALSE, -- if `false` then it's per 100g
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
    CONSTRAINT unit_type_check CHECK (unit_type IN ('items', 'grams', 'ml', 'portions'))
);

-- +goose Down
DROP TABLE recipes;