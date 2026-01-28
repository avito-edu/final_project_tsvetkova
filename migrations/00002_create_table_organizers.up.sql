-- +goose Up
CREATE TABLE organizers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

-- +goose StatementBegin
CREATE INDEX idx_organizers_email ON organizers(email);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS organizers;