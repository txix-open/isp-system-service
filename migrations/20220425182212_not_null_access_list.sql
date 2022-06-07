-- +goose Up
UPDATE access_list SET value = false WHERE value is null;

ALTER TABLE access_list
    ALTER COLUMN value SET NOT NULL;

-- +goose Down
