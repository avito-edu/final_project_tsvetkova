-- +goose Up
DROP TABLE IF EXISTS athletes;

-- +goose Down
CREATE TABLE athletes (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    birth_date DATE,
    status VARCHAR(50),
    gender VARCHAR(10) CHECK (gender IN ('male', 'female'))
);