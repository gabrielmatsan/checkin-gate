package infra

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	identity "github.com/gabrielmatsan/checkin-gate/internal/identity/domain"
	"github.com/jmoiron/sqlx"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *identity.User) (*identity.User, error) {
	query, args, err := psql.
		Insert("users").
		Columns("id", "first_name", "last_name", "email", "role").
		Values(user.ID, user.FirstName, user.LastName, user.Email, user.Role).
		Suffix("RETURNING id, first_name, last_name, email, role, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row identity.User

	err = r.db.GetContext(ctx, &row, query, args...)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*identity.User, error) {
	query, args, err := psql.
		Select("id", "first_name", "last_name", "email", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row identity.User
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*identity.User, error) {
	query, args, err := psql.
		Select("id", "first_name", "last_name", "email", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row identity.User
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *identity.User) error {
	query, args, err := psql.
		Update("users").
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("email", user.Email).
		Set("role", user.Role).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query, args, err := psql.
		Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
