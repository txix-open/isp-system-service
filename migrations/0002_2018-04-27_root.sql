-- +goose Up
INSERT INTO system (id, name) VALUES (1, 'rootSystem');
SELECT setval(concat(current_schema(), '.system_id_seq') :: regclass, 1);

-- +goose Down
DELETE FROM system
WHERE id = 1;