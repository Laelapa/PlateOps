-- +goose Up
CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    private BOOLEAN DEFAULT TRUE,
    description TEXT,
    instructions TEXT,
    unit_type VARCHAR(50) NOT NULL, -- defines the units used for the quantity metric
    quantity INT,
    nutrients_per_portion BOOLEAN DEFAULT FALSE, -- if `false` then it's per 100g
    calories REAL DEFAULT 0 NOT NULL,
    fats REAL DEFAULT 0 NOT NULL,
    saturated REAL DEFAULT 0 NOT NULL,
    carbs REAL DEFAULT 0 NOT NULL,
    sugars REAL DEFAULT 0 NOT NULL,
    protein REAL DEFAULT 0 NOT NULL,
    fiber REAL DEFAULT 0 NOT NULL,
    sodium REAL DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unit_type_check CHECK (unit_type IN ('items', 'grams', 'ml', 'portions'))
);

-- +goose Down
DROP TABLE recipes;