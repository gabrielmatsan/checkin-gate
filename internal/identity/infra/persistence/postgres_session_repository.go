package persistence

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/entity"
	"github.com/jmoiron/sqlx"
)

type PostgresSessionRepository struct {
	db *sqlx.DB
}

func NewPostgresSessionRepository(db *sqlx.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) Save(ctx context.Context, session *entity.Session) error {
	query, args, err := psql.
		Insert("sessions").
		Columns("id", "user_id", "refresh_token", "ip_address", "user_agent", "expires_at", "created_at").
		Values(session.ID, session.UserID, session.RefreshToken, session.IpAddress, session.UserAgent, session.ExpiresAt, session.CreatedAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresSessionRepository) FindByRefreshToken(ctx context.Context, token string) (*entity.Session, error) {
	query, args, err := psql.
		Select("id", "user_id", "refresh_token", "ip_address", "user_agent", "expires_at", "created_at").
		From("sessions").
		Where(sq.Eq{"refresh_token": token}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.Session
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	query, args, err := psql.
		Select("id", "user_id", "refresh_token", "ip_address", "user_agent", "expires_at", "created_at").
		From("sessions").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entity.Session
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.Session, len(rows))
	for i := range rows {
		result[i] = &rows[i]
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
