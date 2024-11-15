-- +goose Up
update service set domain_id = 1, updated_at = now();

-- +goose Down

