package events

import (
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type Activity struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	EventID     string     `json:"event_id" db:"event_id"`
	Description *string    `json:"description" db:"description"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     time.Time  `json:"end_date" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}

type NewActivityParams struct {
	Name        string    `json:"name" db:"name"`
	EventID     string    `json:"event_id" db:"event_id"`
	Description *string   `json:"description" db:"description"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     time.Time `json:"end_date" db:"end_date"`
}

func NewActivity(params NewActivityParams) (*Activity, error) {

	id, err := lib.GenerateID(lib.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate activity ID: %w", err)
	}

	return &Activity{
		ID:          id,
		Name:        params.Name,
		EventID:     params.EventID,
		Description: params.Description,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		CreatedAt:   time.Now(),
	}, nil
}

func (a *Activity) touch() {
	now := time.Now()
	a.UpdatedAt = &now
}

func (a *Activity) Update(params NewActivityParams) error {
	a.Name = params.Name
	a.Description = params.Description
	a.StartDate = params.StartDate
	a.EndDate = params.EndDate
	a.touch()
	return nil
}
