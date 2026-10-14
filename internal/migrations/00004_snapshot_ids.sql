-- +goose Up
ALTER TABLE epochs ADD COLUMN snapshot_id UUID;

-- +goose Down
ALTER TABLE epochs DROP COLUMN snapshot_id;
