package persistence

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/shared"
	"github.com/jmoiron/sqlx"
)

type PostgresCheckInRepository struct {
	db shared.DBTX
}

func NewPostgresCheckInRepository(db shared.DBTX) *PostgresCheckInRepository {
	return &PostgresCheckInRepository{db: db}
}

// WithTx retorna uma nova instância do repositório usando a transação fornecida
func (r *PostgresCheckInRepository) WithTx(tx *sqlx.Tx) *PostgresCheckInRepository {
	return &PostgresCheckInRepository{db: tx}
}

func (r *PostgresCheckInRepository) Save(ctx context.Context, checkIn *entity.CheckIn) (*entity.CheckIn, error) {
	query, args, err := psql.
		Insert("check_ins").
		Columns("id", "user_id", "activity_id", "checked_at").
		Values(checkIn.ID, checkIn.UserID, checkIn.ActivityID, checkIn.CheckedAt).
		Suffix("RETURNING id, user_id, activity_id, checked_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.CheckIn
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresCheckInRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.CheckIn, error) {
	query, args, err := psql.
		Select("id", "user_id", "activity_id", "checked_at").
		From("check_ins").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entity.CheckIn
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.CheckIn, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

func (r *PostgresCheckInRepository) FindByActivityID(ctx context.Context, activityID string) ([]*entity.CheckIn, error) {
	query, args, err := psql.
		Select("id", "user_id", "activity_id", "checked_at").
		From("check_ins").
		Where(sq.Eq{"activity_id": activityID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entity.CheckIn
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.CheckIn, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

func (r *PostgresCheckInRepository) FindByID(ctx context.Context, id string) (*entity.CheckIn, error) {
	query, args, err := psql.
		Select("id", "user_id", "activity_id", "checked_at").
		From("check_ins").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.CheckIn
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresCheckInRepository) FindByUserAndActivity(ctx context.Context, userID, activityID string) (*entity.CheckIn, error) {
	query, args, err := psql.
		Select("id", "user_id", "activity_id", "checked_at").
		From("check_ins").
		Where(sq.Eq{"user_id": userID, "activity_id": activityID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.CheckIn
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresCheckInRepository) FindByActivityIDs(ctx context.Context, activityIDs []string) ([]*entity.CheckIn, error) {
	query, args, err := psql.
		Select("id", "user_id", "activity_id", "checked_at").
		From("check_ins").
		Where(sq.Eq{"activity_id": activityIDs}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var rows []entity.CheckIn
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.CheckIn, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
