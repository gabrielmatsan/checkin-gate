package persistence

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostgresEventRepository struct {
	db *sqlx.DB
}

func NewPostgresEventRepository(db *sqlx.DB) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

func (r *PostgresEventRepository) Save(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	query, args, err := psql.
		Insert("events").
		Columns("id", "name", "allowed_domains", "description", "start_date", "end_date").
		Values(event.ID, event.Name, pq.StringArray(event.AllowedDomains), event.Description, event.StartDate, event.EndDate).
		Suffix("RETURNING id, name, allowed_domains, description, start_date, end_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.Event
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresEventRepository) FindByID(ctx context.Context, id string) (*entity.Event, error) {
	query, args, err := psql.
		Select("id", "name", "allowed_domains", "description", "start_date", "end_date", "created_at", "updated_at").
		From("events").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.Event
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresEventRepository) FindAll(ctx context.Context) ([]*entity.Event, error) {
	query, args, err := psql.
		Select("id", "name", "allowed_domains", "description", "start_date", "end_date", "created_at", "updated_at").
		From("events").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entity.Event
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.Event, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresEventRepository) Update(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	query, args, err := psql.
		Update("events").
		Set("name", event.Name).
		Set("allowed_domains", pq.StringArray(event.AllowedDomains)).
		Set("description", event.Description).
		Set("start_date", event.StartDate).
		Set("end_date", event.EndDate).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": event.ID}).
		Suffix("RETURNING id, name, allowed_domains, description, start_date, end_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.Event
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresEventRepository) Delete(ctx context.Context, id string) error {
	query, args, err := psql.
		Delete("events").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
