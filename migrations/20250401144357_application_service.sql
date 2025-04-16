-- +goose Up
ALTER TABLE application ALTER COLUMN id DROP DEFAULT;

DROP SEQUENCE application_id_seq;

ALTER TABLE service RENAME TO application_group;

ALTER TABLE application
RENAME COLUMN service_id TO application_group_id;

ALTER INDEX uq_name_service_id
RENAME TO uq_name_application_group_id;

ALTER TABLE application
RENAME CONSTRAINT "fk_service_id__domain_id" TO "fk_application_group_id";

-- +goose Down
ALTER TABLE application
RENAME CONSTRAINT "fk_application_group_id" TO "fk_service_id__domain_id";

ALTER INDEX uq_name_application_group_id
RENAME TO uq_name_service_id;

ALTER TABLE application
RENAME COLUMN application_group_id TO service_id;

ALTER TABLE application_group RENAME TO service;

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