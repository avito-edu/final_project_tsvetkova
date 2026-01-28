-- +goose Up
CREATE TABLE competitions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    date DATE NOT NULL,
    organizer_id INTEGER NOT NULL REFERENCES organizers(id) ON DELETE CASCADE,
    results JSONB DEFAULT '[]'
);

-- +goose StatementBegin
CREATE INDEX idx_competitions_date ON competitions(date);
CREATE INDEX idx_competitions_organizer_id ON competitions(organizer_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS competitions;