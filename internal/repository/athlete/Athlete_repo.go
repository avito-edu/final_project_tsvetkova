package athlete

import (
	"database/sql"
	"errors"
	"swim_service/internal/domain"

	"github.com/Masterminds/squirrel"
)

type AthleteRepository struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

func NewAthleteRepository(db *sql.DB) *AthleteRepository {
	return &AthleteRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *AthleteRepository) CreateAthlete(athlete *domain.Athlete) error {
	query, args, err := r.psql.
		Insert("athletes").
		Columns("name", "birth_date", "status", "gender").
		Values(athlete.Name, athlete.BirthDate, athlete.Status, athlete.Gender).
		Suffix("RETURNING user_id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&athlete.ID)
	return err
}

func (r *AthleteRepository) GetByID(id int) (*domain.Athlete, error) {
	query, args, err := r.psql.
		Select("user_id", "name", "birth_date", "status", "gender").
		From("athletes").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	athlete := &domain.Athlete{}
	err = r.db.QueryRow(query, args...).Scan(
		&athlete.ID,
		&athlete.Name,
		&athlete.BirthDate,
		&athlete.Status,
		&athlete.Gender,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return athlete, err
}

func (r *AthleteRepository) UpdateStatus(id int, newStatus string) error {
	query, args, err := r.psql.
		Update("athletes").
		Set("status", newStatus).
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)
	return err
}
