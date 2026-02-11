package infra

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	identity "github.com/gabrielmatsan/checkin-gate/internal/identity/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresSessionRepository struct {
	db *sqlx.DB
}

func NewPostgresSessionRepository(db *sqlx.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) Save(ctx context.Context, session *identity.Session) error {
	query, args, err := psql.
		Insert("sessions").
		Columns("id", "user_id", "refresh_token", "ip_address", "user_agent", "expires_at", "created_at").
		Values(
			session.ID,
			session.UserID,
			session.RefreshToken,
			session.IpAddress,
			session.UserAgent,
			session.ExpiresAt,
			session.CreatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresSessionRepository) FindByRefreshToken(ctx context.Context, token string) (*identity.Session, error) {
	query, args, err := psql.
		Select("*").
		From("sessions").
		Where(sq.Eq{"refresh_token": token}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var session identity.Session
	if err := r.db.GetContext(ctx, &session, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (r *PostgresSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*identity.Session, error) {
	query, args, err := psql.
		Select("*").
		From("sessions").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var sessions []identity.Session
	if err := r.db.SelectContext(ctx, &sessions, query, args...); err != nil {
		return nil, err
	}

	result := make([]*identity.Session, len(sessions))
	for i := range sessions {
		result[i] = &sessions[i]
	}

	return result, nil
}

func (r *PostgresSessionRepository) Delete(ctx context.Context, id string) error {
	query, args, err := psql.
		Delete("sessions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresSessionRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	query, args, err := psql.
		Delete("sessions").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresSessionRepository) DeleteExpired(ctx context.Context) error {
	query, args, err := psql.
		Delete("sessions").
		Where(sq.Lt{"expires_at": sq.Expr("NOW()")}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
