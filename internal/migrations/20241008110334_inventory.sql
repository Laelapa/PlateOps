-- +goose Up
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id),
    product_id INT NOT NULL REFERENCES food_registry(product_id),
    expiration_date DATE,
    is_opened BOOLEAN DEFAULT FALSE,
    max_quantity FLOAT NOT NULL,
    curr_quantity FLOAT NOT NULL,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    opened_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE inventory;