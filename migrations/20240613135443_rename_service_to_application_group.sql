-- +goose Up
ALTER TABLE service RENAME TO application_group;

ALTER TABLE application RENAME COLUMN service_id TO application_group_id;

ALTER TABLE application RENAME CONSTRAINT UQ_name_service_id TO UQ_name_application_group_id;

ALTER TABLE application RENAME CONSTRAINT FK_service_id__domain_id TO FK_application_group_id__domain_id;

ALTER TRIGGER modify_service ON application_group RENAME TO modify_application_group;
-- +goose Down
