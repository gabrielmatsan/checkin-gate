package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/shared"
	"github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostgresEventRepository struct {
	db shared.DBTX
}

func NewPostgresEventRepository(db shared.DBTX) *PostgresEventRepository {
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

func (r *PostgresEventRepository) PartialUpdate(ctx context.Context, id string, input repository.UpdateEventInput) (*entity.Event, error) {
	builder := psql.Update("events").Where(sq.Eq{"id": id})

	if input.Name != nil {
		builder = builder.Set("name", *input.Name)
	}
	if input.AllowedDomains != nil {
		builder = builder.Set("allowed_domains", pq.StringArray(*input.AllowedDomains))
	}
	if input.Description != nil {
		builder = builder.Set("description", *input.Description)
	}
	if input.StartDate != nil {
		builder = builder.Set("start_date", *input.StartDate)
	}
	if input.EndDate != nil {
		builder = builder.Set("end_date", *input.EndDate)
	}
	if input.Status != nil {
		builder = builder.Set("status", *input.Status)
	}

	builder = builder.Set("updated_at", sq.Expr("NOW()"))
	builder = builder.Suffix("RETURNING id, name, allowed_domains, description, start_date, end_date, status, created_at, updated_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.Event
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresEventRepository) FindByIDWithActivitiesAndCheckIns(ctx context.Context, eventID string) (*repository.EventWithActivitiesAndCheckIns, error) {
	var row struct {
		entity.Event
		Activities json.RawMessage `db:"activities"`
	}

	query := `
		SELECT
			e.*,
			COALESCE(
				json_agg(
					json_build_object(
						'activity_id', a.id,
						'activity_name', a.name,
						'check_ins', (
							SELECT COALESCE(json_agg(
								json_build_object(
									'id', c.id,
									'user_id', c.user_id,
									'activity_id', c.activity_id,
									'checked_at', c.checked_at
								)
							), '[]'::json)
							FROM check_ins c
							WHERE c.activity_id = a.id
						)
					)
				) FILTER (WHERE a.id IS NOT NULL),
				'[]'::json
			) AS activities
		FROM events e
		LEFT JOIN activities a ON a.event_id = e.id
		WHERE e.id = $1
		GROUP BY e.id
	`

	if err := r.db.GetContext(ctx, &row, query, eventID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var activities []repository.ActivityWithCheckIns
	if err := json.Unmarshal(row.Activities, &activities); err != nil {
		return nil, err
	}

	return &repository.EventWithActivitiesAndCheckIns{
		Event:      &row.Event,
		Activities: activities,
	}, nil
}
