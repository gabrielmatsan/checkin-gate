package infra

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	events "github.com/gabrielmatsan/checkin-gate/internal/events/domain"
	"github.com/jmoiron/sqlx"
)

type PostgresActivityRepository struct {
	db *sqlx.DB
}

func NewPostgresActivityRepository(db *sqlx.DB) *PostgresActivityRepository {
	return &PostgresActivityRepository{db: db}
}

func (r *PostgresActivityRepository) Save(ctx context.Context, activity *events.Activity) (*events.Activity, error) {
	query, args, err := psql.
		Insert("activities").
		Columns("id", "name", "event_id", "description", "start_date", "end_date").
		Values(activity.ID, activity.Name, activity.EventID, activity.Description, activity.StartDate, activity.EndDate).
		Suffix("RETURNING id, name, event_id, description, start_date, end_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row events.Activity
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresActivityRepository) SaveAll(ctx context.Context, activities []*events.Activity) ([]*events.Activity, error) {
	if len(activities) == 0 {
		return nil, nil
	}

	builder := psql.
		Insert("activities").
		Columns("id", "name", "event_id", "description", "start_date", "end_date")

	for _, a := range activities {
		builder = builder.Values(a.ID, a.Name, a.EventID, a.Description, a.StartDate, a.EndDate)
	}

	query, args, err := builder.
		Suffix("RETURNING id, name, event_id, description, start_date, end_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []events.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*events.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) FindByEventIDAndNames(ctx context.Context, eventID string, names []string) ([]*events.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		Where(sq.Eq{"event_id": eventID, "name": names}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []events.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*events.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) FindByID(ctx context.Context, id string) (*events.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row events.Activity
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresActivityRepository) FindByEventID(ctx context.Context, eventID string) ([]*events.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		Where(sq.Eq{"event_id": eventID}).
		OrderBy("start_date ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []events.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*events.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) FindAll(ctx context.Context) ([]*events.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		OrderBy("start_date ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []events.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*events.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) Update(ctx context.Context, activity *events.Activity) (*events.Activity, error) {
	query, args, err := psql.
		Update("activities").
		Set("name", activity.Name).
		Set("description", activity.Description).
		Set("start_date", activity.StartDate).
		Set("end_date", activity.EndDate).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": activity.ID}).
		Suffix("RETURNING id, name, event_id, description, start_date, end_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row events.Activity
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresActivityRepository) Delete(ctx context.Context, id string) error {
	query, args, err := psql.
		Delete("activities").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
