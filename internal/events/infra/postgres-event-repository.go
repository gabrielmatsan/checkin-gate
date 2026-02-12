package infra

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	events "github.com/gabrielmatsan/checkin-gate/internal/events/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// EventRow is used for scanning PostgreSQL arrays
type EventRow struct {
	ID             string         `db:"id"`
	Name           string         `db:"name"`
	AllowedDomains pq.StringArray `db:"allowed_domains"`
	Description    *string        `db:"description"`
	StartDate      sql.NullTime   `db:"start_date"`
	EndDate        sql.NullTime   `db:"end_date"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
}

func (r *EventRow) ToEvent() *events.Event {
	event := &events.Event{
		ID:             r.ID,
		Name:           r.Name,
		AllowedDomains: r.AllowedDomains,
		Description:    r.Description,
	}

	if r.StartDate.Valid {
		event.StartDate = r.StartDate.Time
	}
	if r.EndDate.Valid {
		event.EndDate = r.EndDate.Time
	}
	if r.CreatedAt.Valid {
		event.CreatedAt = r.CreatedAt.Time
	}
	if r.UpdatedAt.Valid {
		event.UpdatedAt = &r.UpdatedAt.Time
	}

	return event
}

type PostgresEventRepository struct {
	db *sqlx.DB
}

func NewPostgresEventRepository(db *sqlx.DB) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

func (r *PostgresEventRepository) Save(ctx context.Context, event *events.Event) (*events.Event, error) {
	query, args, err := psql.
		Insert("events").
		Columns("id", "name", "allowed_domains", "description", "start_date", "end_date").
		Values(event.ID, event.Name, pq.StringArray(event.AllowedDomains), event.Description, event.StartDate, event.EndDate).
		Suffix("RETURNING id, name, allowed_domains, description, start_date, end_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row EventRow
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return row.ToEvent(), nil
}

func (r *PostgresEventRepository) FindByID(ctx context.Context, id string) (*events.Event, error) {
	query, args, err := psql.
		Select("id", "name", "allowed_domains", "description", "start_date", "end_date", "created_at", "updated_at").
		From("events").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row EventRow
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return row.ToEvent(), nil
}

func (r *PostgresEventRepository) FindAll(ctx context.Context) ([]*events.Event, error) {
	query, args, err := psql.
		Select("id", "name", "allowed_domains", "description", "start_date", "end_date", "created_at", "updated_at").
		From("events").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []EventRow
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*events.Event, len(rows))
	for i, row := range rows {
		result[i] = row.ToEvent()
	}

	return result, nil
}

func (r *PostgresEventRepository) Update(ctx context.Context, event *events.Event) (*events.Event, error) {
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

	var row EventRow
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return row.ToEvent(), nil
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
