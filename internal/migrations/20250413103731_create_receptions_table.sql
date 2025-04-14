-- +goose Up
CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID NOT NULL REFERENCES pvz_points(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_created_at ON receptions(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_created_at;
DROP TABLE IF EXISTS receptions;