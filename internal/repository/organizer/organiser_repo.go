package organizer

import (
	"database/sql"
	"errors"
	"swim_service/internal/domain"

	"github.com/Masterminds/squirrel"
)

type OrganizerRepository struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

func NewOrganizerRepository(db *sql.DB) *OrganizerRepository {
	return &OrganizerRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *OrganizerRepository) Create(organizer *domain.Organizer) error {
	query, args, err := r.psql.
		Insert("organizers").
		Columns("name", "email").
		Values(organizer.Name, organizer.Email).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(query, args...).Scan(&organizer.ID)
	return err
}

func (r *OrganizerRepository) GetByID(id int) (*domain.Organizer, error) {
	query, args, err := r.psql.
		Select("id", "name", "email").
		From("organizers").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	organizer := &domain.Organizer{}
	err = r.db.QueryRow(query, args...).Scan(&organizer.ID, &organizer.Name, &organizer.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return organizer, err
}
