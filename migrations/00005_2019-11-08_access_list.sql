-- +goose Up
CREATE TABLE access_list
(
    app_id INT,
    method VARCHAR(255),
    value  BOOLEAN,
    PRIMARY KEY (app_id, method)
);
CREATE INDEX IX_access_list__app_id ON access_list (app_id);
-- +goose Down