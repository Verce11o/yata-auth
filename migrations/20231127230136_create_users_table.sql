-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY ,
    username VARCHAR(255) NOT NULL UNIQUE ,
    email VARCHAR(255) NOT NULL UNIQUE ,
    password VARChAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
