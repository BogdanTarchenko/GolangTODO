-- +goose Up
CREATE TABLE tasks
(
    id           VARCHAR PRIMARY KEY,
    title        VARCHAR   NOT NULL,
    description  TEXT,
    deadline     TIMESTAMP,
    status       VARCHAR   NOT NULL,
    priority     VARCHAR   NOT NULL,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP,
    is_completed BOOLEAN   NOT NULL
);

-- +goose Down
DROP TABLE tasks;