package entity

import (
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type CheckIn struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	ActivityID string    `db:"activity_id"`
	CheckedAt  time.Time `db:"checked_at"`
}

type NewCheckInParams struct {
	UserID     string
	ActivityID string
}

func NewCheckIn(params NewCheckInParams) (*CheckIn, error) {
	id, err := lib.GenerateID(lib.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate check-in ID: %w", err)
	}

	return &CheckIn{
		ID:         id,
		UserID:     params.UserID,
		ActivityID: params.ActivityID,
		CheckedAt:  time.Now(),
	}, nil
}
