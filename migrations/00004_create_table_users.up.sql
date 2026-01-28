-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL
);

-- +goose StatementBegin
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_role ON users(role);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS users;