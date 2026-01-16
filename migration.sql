CREATE TABLE IF NOT EXISTS athletes (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    birth_date DATE,
    status VARCHAR(50),
    gender VARCHAR(10) CHECK (gender IN ('male', 'female'))
);
CREATE TABLE IF NOT EXISTS organizers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS competitions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    date DATE NOT NULL,
    organizer_id INTEGER NOT NULL REFERENCES organizers(id) ON DELETE CASCADE,
    results JSONB DEFAULT '[]'
);
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL
);