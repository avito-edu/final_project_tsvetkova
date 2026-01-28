-- +goose Up
CREATE TABLE athletes (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    birth_date DATE,
    status VARCHAR(50),
    gender VARCHAR(10) CHECK (gender IN ('male', 'female'))
);

-- +goose StatementBegin
CREATE INDEX idx_athletes_name ON athletes(name);
CREATE INDEX idx_athletes_status ON athletes(status);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS athletes;