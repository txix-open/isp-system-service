-- +goose Up
CREATE TABLE access_list
(
    app_id INT,
    method VARCHAR(255),
    value  BOOLEAN,
    PRIMARY KEY (app_id, method),
    CONSTRAINT FK_app_id__application_id FOREIGN KEY (app_id) REFERENCES application (id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose Down