-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_food_registry_name_trgm
ON food_registry USING gin (name gin_trgm_ops);

-- +goose Down
DROP INDEX IF EXISTS idx_food_registry_name_trgm;
