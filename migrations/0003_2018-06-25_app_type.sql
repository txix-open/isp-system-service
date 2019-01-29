-- +goose Up
ALTER TABLE application
  ADD COLUMN type VARCHAR(255) NOT NULL DEFAULT 'SYSTEM';

-- +goose Down
ALTER TABLE application
  DROP COLUMN "type";