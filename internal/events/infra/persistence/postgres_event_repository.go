package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
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
		WHERE e.id = :event_id
		GROUP BY e.id
	`

	nstmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer nstmt.Close()

	if err := nstmt.GetContext(ctx, &row, map[string]interface{}{"event_id": eventID}); err != nil {
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
