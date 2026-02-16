package entity

import (
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type Activity struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	EventID     string     `db:"event_id"`
	Description *string    `db:"description"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     time.Time  `db:"end_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type NewActivityParams struct {
	Name        string
	EventID     string
	Description *string
	StartDate   time.Time
	EndDate     time.Time
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
		UpdatedAt:   nil,
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

// starttime antes de endtime
func (a *Activity) IsStartDateBeforeEndDate() bool {
	return a.StartDate.Before(a.EndDate)
}

// endtime depois de starttime
func (a *Activity) IsEndDateAfterStartDate() bool {
	return a.EndDate.After(a.StartDate)
}

func (a *Activity) HasStarted() bool {
	return a.StartDate.Before(time.Now())
}

func (a *Activity) HasEnded() bool {
	return a.EndDate.Before(time.Now())
}

func (a *Activity) IsCheckInAllowed(checkInTime time.Time) bool {
	return a.HasStarted() && !a.HasEnded()
}
