package events

import (
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type Event struct {
	ID             string     `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	AllowedDomains []string   `json:"allowed_domains" db:"allowed_domains"`
	Description    *string    `json:"description" db:"description"`
	StartDate      time.Time  `json:"start_date" db:"start_date"`
	EndDate        time.Time  `json:"end_date" db:"end_date"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}

type NewEventParams struct {
	Name           string    `json:"name"`
	AllowedDomains *[]string `json:"allowed_domains"`
	Description    *string   `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

func NewEvent(params NewEventParams) (*Event, error) {
	id, err := lib.GenerateID(lib.UUID)

	if err != nil {
		return nil, fmt.Errorf("failed to generate event ID: %w", err)
	}

	return &Event{
		ID:             id,
		Name:           params.Name,
		AllowedDomains: *params.AllowedDomains,
		Description:    params.Description,
		StartDate:      params.StartDate,
		EndDate:        params.EndDate,
		CreatedAt:      time.Now(),
	}, nil
}
