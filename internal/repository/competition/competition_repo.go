package competition

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"swim_service/internal/domain"
	"time"

	"github.com/Masterminds/squirrel"
)

type PostgresCompetitionRepository struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

func NewPostgresCompetitionRepository(db *sql.DB) *PostgresCompetitionRepository {
	return &PostgresCompetitionRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresCompetitionRepository) Create(comp *domain.Competition) error {
	resultsJSON, _ := json.Marshal(comp.Results)

	query, args, err := r.psql.
		Insert("competitions").
		Columns("name", "date", "organizer_id", "results").
		Values(comp.Name, comp.Date, comp.Organizer, resultsJSON).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&comp.ID)
	return err
}

func (r *PostgresCompetitionRepository) GetByID(id int) (*domain.Competition, error) {
	query, args, err := r.psql.
		Select("id", "name", "date", "organizer_id", "results").
		From("competitions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var comp domain.Competition
	var resultsJSON []byte

	err = r.db.QueryRow(query, args...).Scan(&comp.ID, &comp.Name, &comp.Date, &comp.Organizer, &resultsJSON)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resultsJSON, &comp.Results); err != nil {
		log.Printf("Failed to unmarshal: %s", err.Error())
		return nil, err
	}

	return &comp, nil
}

func (r *PostgresCompetitionRepository) GetAll() ([]*domain.Competition, error) {
	query, args, err := r.psql.
		Select("id", "name", "date", "organizer_id", "results").
		From("competitions").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competitions []*domain.Competition
	for rows.Next() {
		var comp domain.Competition
		var resultsJSON []byte

		err := rows.Scan(&comp.ID, &comp.Name, &comp.Date, &comp.Organizer, &resultsJSON)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(resultsJSON, &comp.Results); err != nil {
			log.Printf("Failed to unmarshal: %s", err.Error())
			continue
		}

		competitions = append(competitions, &comp)
	}

	return competitions, nil
}

func (r *PostgresCompetitionRepository) AddResult(compID int, result domain.Result) error {
	query, args, err := r.psql.
		Select("results").
		From("competitions").
		Where(squirrel.Eq{"id": compID}).
		ToSql()
	if err != nil {
		return err
	}

	var resultsJSON []byte
	var currentResults []domain.Result

	err = r.db.QueryRow(query, args...).Scan(&resultsJSON)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resultsJSON, &currentResults); err != nil {
		return err
	}

	currentResults = append(currentResults, result)
	updatedJSON, _ := json.Marshal(currentResults)

	updateQuery, updateArgs, err := r.psql.
		Update("competitions").
		Set("results", updatedJSON).
		Where(squirrel.Eq{"id": compID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(updateQuery, updateArgs...)
	return err
}

func (r *PostgresCompetitionRepository) GetAthleteResultsByDistance(athleteID int, distance string) ([]domain.Result, error) {
	query := `
		SELECT jsonb_array_elements(results) as result
		FROM competitions
		WHERE results @> $1::jsonb
	`
	filter := []map[string]any{
		{"athlete_id": athleteID, "distance": distance},
	}
	filterJSON, _ := json.Marshal(filter)

	rows, err := r.db.Query(query, filterJSON)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Result
	for rows.Next() {
		var resultJSON []byte
		var result domain.Result

		err := rows.Scan(&resultJSON)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(resultJSON, &result); err != nil {
			log.Printf("Failed to unmarshal: %s", err.Error())
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *PostgresCompetitionRepository) GetAllResultsByDistance(distance string) ([]domain.Result, error) {
	query := `
		SELECT jsonb_array_elements(results) as result
		FROM competitions
		WHERE results @> $1::jsonb
	`

	filter := []map[string]any{
		{"distance": distance},
	}
	filterJSON, _ := json.Marshal(filter)

	rows, err := r.db.Query(query, filterJSON)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Result
	for rows.Next() {
		var resultJSON []byte
		var result domain.Result

		err := rows.Scan(&resultJSON)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(resultJSON, &result); err != nil {
			log.Printf("Failed to unmarshal: %s", err.Error())
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *PostgresCompetitionRepository) GetAthleteResultsByPeriod(ctx context.Context, athleteID int, startDate, endDate time.Time) ([]domain.Result, error) {
	query := `
		SELECT 
			(r.elem->>'AthleteID')::int as athlete_id,
			r.elem->>'Distance' as distance,
			(r.elem->>'TimeSec')::float as time_sec,
			(r.elem->>'Place')::int as place
		FROM competitions c,
		LATERAL jsonb_array_elements(c.results) AS r(elem)
		WHERE c.date BETWEEN $1 AND $2
		AND (r.elem->>'AthleteID')::int = $3;
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, athleteID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []domain.Result
	for rows.Next() {
		var result domain.Result
		err := rows.Scan(
			&result.AthleteID,
			&result.Distance,
			&result.TimeSec,
			&result.Place,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}
