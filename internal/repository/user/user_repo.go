package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"swim_service/internal/domain"

	"github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *sql.DB
	sq squirrel.StatementBuilderType
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query, args, err := r.sq.
		Insert("users").
		Columns("login", "password", "role").
		Values(user.Login, user.Password, user.Role).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil
}

func (r *UserRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	query, args, err := r.sq.
		Select("1").
		From("users").
		Where(squirrel.Eq{"login": login}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build exists query: %w", err)
	}

	var dummy int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&dummy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return true, nil
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	query, args, err := r.sq.
		Select("id", "login", "password", "role").
		From("users").
		Where(squirrel.Eq{"login": login}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build find query: %w", err)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	user := &domain.User{}
	err = row.Scan(&user.ID, &user.Login, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return user, nil
}
