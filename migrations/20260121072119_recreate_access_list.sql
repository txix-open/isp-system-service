-- +goose Up
CREATE TABLE access_list_new
(
    app_id INT,
    http_method VARCHAR(255) NOT NULL DEFAULT '',
    method VARCHAR(255),
    value  BOOLEAN NOT NULL,
    PRIMARY KEY (app_id, http_method, method),
    CONSTRAINT FK_app_id__application_id FOREIGN KEY (app_id) REFERENCES application (id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO access_list_new (app_id, method, value)
SELECT app_id, method, value
FROM access_list
WHERE value = true;

DROP TABLE access_list;
ALTER TABLE access_list_new RENAME TO access_list;

-- +goose Down
CREATE TABLE access_list_old
(
    app_id INT,
    method VARCHAR(255),
    value  BOOLEAN NOT NULL,
    PRIMARY KEY (app_id, method),
    CONSTRAINT FK_app_id__application_id FOREIGN KEY (app_id) REFERENCES application (id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO access_list_old (app_id, method, value)
SELECT app_id, method, value
FROM access_list
WHERE value = true;

DROP TABLE access_list;
ALTER TABLE access_list_old RENAME TO access_list;