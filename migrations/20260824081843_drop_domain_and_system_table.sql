-- +goose Up
ALTER TABLE application_group
    DROP CONSTRAINT fk_domain_id__domain_id;
ALTER TABLE application_group
    DROP CONSTRAINT uq_name_domain_name;

ALTER TABLE application_group
    DROP COLUMN domain_id;

ALTER TABLE application_group
    ADD CONSTRAINT uq_application_group_name UNIQUE ("name");

ALTER TABLE domain
    DROP CONSTRAINT fk_system_id__system_id;

ALTER TABLE domain
    DROP CONSTRAINT uq_name_system_id;

DROP TABLE domain;
DROP TABLE system;

-- +goose Down

-- irreversible migration