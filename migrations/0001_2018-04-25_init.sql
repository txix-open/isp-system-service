-- +goose Up
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
  FOREIGN KEY (system_id) REFERENCES system (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT UQ_name_system_id UNIQUE (name, system_id)
);

CREATE TABLE service (
  id          SERIAL4      NOT NULL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  description TEXT,
  domain_id   INTEGER      NOT NULL,
  created_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  CONSTRAINT FK_domain_id__domain_id
  FOREIGN KEY (domain_id) REFERENCES domain (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT UQ_name_domain_name UNIQUE (name, domain_id)
);

CREATE TABLE application (
  id          SERIAL4      NOT NULL PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  description TEXT,
  service_id  INTEGER      NOT NULL,
  created_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at  TIMESTAMP    NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  CONSTRAINT FK_service_id__domain_id
  FOREIGN KEY (service_id) REFERENCES service (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT UQ_name_service_id UNIQUE (name, service_id)
);

CREATE TABLE token (
  token       TEXT      NOT NULL PRIMARY KEY,
  created_at  TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
  app_id      INTEGER   NOT NULL,
  expire_time BIGINT    NOT NULL DEFAULT -1,
  CONSTRAINT FK_app_id__application_id
  FOREIGN KEY (app_id) REFERENCES application (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_created_modified_column_date()
  RETURNS TRIGGER AS
$body$
BEGIN
  IF TG_OP = 'UPDATE'
  THEN
    NEW.created_at = OLD.created_at;
    NEW.updated_at = (now() at time zone 'utc');
  ELSIF TG_OP = 'INSERT'
    THEN
      NEW.updated_at = (now() at time zone 'utc');
  END IF;
  RETURN NEW;
END;
$body$
LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER modify_system
  BEFORE UPDATE OR INSERT
  ON system
  FOR EACH ROW EXECUTE PROCEDURE update_created_modified_column_date();

CREATE TRIGGER modify_domain
  BEFORE UPDATE OR INSERT
  ON domain
  FOR EACH ROW EXECUTE PROCEDURE update_created_modified_column_date();

CREATE TRIGGER modify_service
  BEFORE UPDATE OR INSERT
  ON service
  FOR EACH ROW EXECUTE PROCEDURE update_created_modified_column_date();

CREATE TRIGGER modify_application
  BEFORE UPDATE OR INSERT
  ON application
  FOR EACH ROW EXECUTE PROCEDURE update_created_modified_column_date();

-- +goose Down
DROP SCHEMA IF EXISTS system_service CASCADE;