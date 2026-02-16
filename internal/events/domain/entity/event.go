package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/lib/pq"
)

// Enum status do evento
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusCompleted EventStatus = "completed"
)

type Event struct {
	ID             string         `db:"id"`
	Name           string         `db:"name"`
	AllowedDomains pq.StringArray `db:"allowed_domains"`
	Description    *string        `db:"description"`
	StartDate      time.Time      `db:"start_date"`
	EndDate        time.Time      `db:"end_date"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      *time.Time     `db:"updated_at"`
	Status         EventStatus    `db:"status"`
}

type NewEventParams struct {
	Name           string
	AllowedDomains *[]string
	Description    *string
	StartDate      time.Time
	EndDate        time.Time
}

func NewEvent(params NewEventParams) (*Event, error) {
	id, err := lib.GenerateID(lib.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate event ID: %w", err)
	}

	var domains pq.StringArray
	if params.AllowedDomains != nil {
		domains = *params.AllowedDomains
	}

	return &Event{
		ID:             id,
		Name:           params.Name,
		AllowedDomains: domains,
		Description:    params.Description,
		StartDate:      params.StartDate,
		EndDate:        params.EndDate,
		CreatedAt:      time.Now(),
		UpdatedAt:      nil,
		Status:         EventStatusDraft,
	}, nil
}

// func (e *Event) touch() {
// 	now := time.Now()
// 	e.UpdatedAt = &now
// }

// verifica se o dominio passado é valido
// se AllowedDomains for nulo ou vazio, permite todos os domínios
func (e *Event) IsAllowedDomain(email string) bool {
	if len(e.AllowedDomains) == 0 {
		return true
	}

	domain := extractDomain(email)
	for _, allowedDomain := range e.AllowedDomains {
		if allowedDomain == domain {
			return true
		}
	}
	return false
}

// verifica se o check-in está dentro do horário de evento
func (e *Event) IsCheckInWithinEventTime(checkInTime time.Time) bool {
	return checkInTime.After(e.StartDate) && checkInTime.Before(e.EndDate)
}

// verifica se a data de inicio é antes da data de fim
func (e *Event) IsStartDateBeforeEndDate() bool {
	return e.StartDate.Before(e.EndDate)
}

// verifica se a data de fim é depois da data de inicio
func (e *Event) IsEndDateAfterStartDate() bool {
	return e.EndDate.After(e.StartDate)
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
