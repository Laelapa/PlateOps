-- +goose Up
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES food_registry(product_id) ON DELETE CASCADE,
    expiration_date DATE,
    is_opened BOOLEAN DEFAULT FALSE,
    current_quantity FLOAT NOT NULL,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    opened_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE inventory;