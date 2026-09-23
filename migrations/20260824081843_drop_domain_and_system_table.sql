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
CREATE TABLE system (
    id          SERIAL4      NOT NULL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    CONSTRAINT UQ_name UNIQUE (name)
);

CREATE TABLE domain (
    id          SERIAL4      NOT NULL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    system_id   INTEGER      NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    CONSTRAINT FK_system_id__system_id
        FOREIGN KEY (system_id)
        REFERENCES system (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT UQ_name_system_id UNIQUE (name, system_id)
);

-- Восстанавливаем root system/domain.
INSERT INTO system (id, name)
VALUES (1, 'rootSystem');

SELECT setval(
    'system_id_seq',
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM system), 1)
);

INSERT INTO domain (id, name, system_id)
VALUES (1, 'root', 1);

SELECT setval(
    'domain_id_seq',
    GREATEST((SELECT COALESCE(MAX(id), 1) FROM domain), 1)
);

ALTER TABLE application_group
    DROP CONSTRAINT uq_application_group_name;

ALTER TABLE application_group
    ADD COLUMN domain_id INTEGER DEFAULT 1;

ALTER TABLE application_group
    ALTER COLUMN domain_id SET NOT NULL;

ALTER TABLE application_group
    ADD CONSTRAINT FK_domain_id__domain_id
        FOREIGN KEY (domain_id)
        REFERENCES domain (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE;

ALTER TABLE application_group
    ADD CONSTRAINT UQ_name_domain_name
        UNIQUE (name, domain_id);

CREATE TRIGGER modify_system
    BEFORE UPDATE OR INSERT ON system
    FOR EACH ROW EXECUTE PROCEDURE update_created_modified_column_date();

CREATE TRIGGER modify_domain
    BEFORE UPDATE OR INSERT ON domain
    FOR EACH ROW EXECUTE PROCEDURE update_created_modified_column_date();