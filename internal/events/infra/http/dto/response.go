package dto

import "time"

type EventResponse struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	AllowedDomains []string   `json:"allowed_domains"`
	Description    *string    `json:"description,omitempty"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type ActivityResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	EventID     string     `json:"event_id"`
	Description *string    `json:"description,omitempty"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type EventWithActivitiesResponse struct {
	Event      EventResponse      `json:"event"`
	Activities []ActivityResponse `json:"activities"`
}

type EventListResponse struct {
	Events []EventResponse `json:"events"`
	Total  int             `json:"total"`
}

type ActivitiesResponse struct {
	Activities []ActivityResponse `json:"activities"`
}
