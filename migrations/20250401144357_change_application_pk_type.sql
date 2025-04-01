-- +goose Up
ALTER TABLE application
ALTER COLUMN id
DROP DEFAULT;

DROP SEQUENCE application_id_seq;

-- +goose Down
CREATE SEQUENCE application_id_seq;

SELECT setval(
        'application_id_seq', COALESCE(
            (
                SELECT MAX(id)
                FROM isp_system_service.application
            ), 1
        )
    );

ALTER TABLE application
ALTER COLUMN id
SET DEFAULT nextval('application_id_seq');