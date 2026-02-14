package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresActivityRepository struct {
	db *sqlx.DB
}

func NewPostgresActivityRepository(db *sqlx.DB) *PostgresActivityRepository {
	return &PostgresActivityRepository{db: db}
}

func (r *PostgresActivityRepository) Save(ctx context.Context, activity *entity.Activity) (*entity.Activity, error) {
	query, args, err := psql.
		Insert("activities").
		Columns("id", "name", "event_id", "description", "start_date", "end_date").
		Values(activity.ID, activity.Name, activity.EventID, activity.Description, activity.StartDate, activity.EndDate).
		Suffix("RETURNING id, name, event_id, description, start_date, end_date, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.Activity
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *PostgresActivityRepository) SaveAll(ctx context.Context, activities []*entity.Activity) ([]*entity.Activity, error) {
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

	var rows []entity.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) FindByEventIDAndNames(ctx context.Context, eventID string, names []string) ([]*entity.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		Where(sq.Eq{"event_id": eventID, "name": names}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entity.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) FindByID(ctx context.Context, id string) (*entity.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var row entity.Activity
	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *PostgresActivityRepository) FindByEventID(ctx context.Context, eventID string) ([]*entity.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		Where(sq.Eq{"event_id": eventID}).
		OrderBy("start_date ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entity.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) FindAll(ctx context.Context) ([]*entity.Activity, error) {
	query, args, err := psql.
		Select("id", "name", "event_id", "description", "start_date", "end_date", "created_at", "updated_at").
		From("activities").
		OrderBy("start_date ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	var rows []entity.Activity
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	result := make([]*entity.Activity, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}

	return result, nil
}

func (r *PostgresActivityRepository) Update(ctx context.Context, activity *entity.Activity) (*entity.Activity, error) {
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

	var row entity.Activity
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


func (r *PostgresActivityRepository) FindByActivityIDWithEvent(ctx context.Context, activityID string) (*repository.ActivityWithEvent, error) {
	// Struct para scan do JOIN
	var row struct {
		// Activity fields
		ID          string     `db:"id"`
		Name        string     `db:"name"`
		EventID     string     `db:"event_id"`
		Description *string    `db:"description"`
		StartDate   time.Time  `db:"start_date"`
		EndDate     time.Time  `db:"end_date"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at"`
		// Event fields
		EventName           string         `db:"event_name"`
		EventAllowedDomains pq.StringArray `db:"event_allowed_domains"`
		EventDescription    *string        `db:"event_description"`
		EventStartDate      time.Time      `db:"event_start_date"`
		EventEndDate        time.Time      `db:"event_end_date"`
		EventCreatedAt      time.Time      `db:"event_created_at"`
		EventUpdatedAt      *time.Time     `db:"event_updated_at"`
	}

	query, args, err := psql.
		Select(
			"a.id", "a.name", "a.event_id", "a.description", "a.start_date", "a.end_date", "a.created_at", "a.updated_at",
			"e.name AS event_name", "e.allowed_domains AS event_allowed_domains", "e.description AS event_description",
			"e.start_date AS event_start_date", "e.end_date AS event_end_date", "e.created_at AS event_created_at", "e.updated_at AS event_updated_at",
		).
		From("activities a").
		InnerJoin("events e ON a.event_id = e.id").
		Where(sq.Eq{"a.id": activityID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &repository.ActivityWithEvent{
		Activity: &entity.Activity{
			ID:          row.ID,
			Name:        row.Name,
			EventID:     row.EventID,
			Description: row.Description,
			StartDate:   row.StartDate,
			EndDate:     row.EndDate,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		},
		Event: &entity.Event{
			ID:             row.EventID,
			Name:           row.EventName,
			AllowedDomains: row.EventAllowedDomains,
			Description:    row.EventDescription,
			StartDate:      row.EventStartDate,
			EndDate:        row.EventEndDate,
			CreatedAt:      row.EventCreatedAt,
			UpdatedAt:      row.EventUpdatedAt,
		},
	}, nil
} 