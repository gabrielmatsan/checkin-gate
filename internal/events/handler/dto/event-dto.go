package dto

import "time"

// CreateEventRequest represents the request body for creating an event
type CreateEventRequest struct {
	Name           string    `json:"name" validate:"required,min=3,max=255"`
	AllowedDomains []string  `json:"allowed_domains" validate:"dive,fqdn"`
	Description    *string   `json:"description" validate:"omitempty,max=500"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

// UpdateEventRequest represents the request body for updating an event
type UpdateEventRequest struct {
	Name           *string    `json:"name" validate:"omitempty,min=3,max=255"`
	AllowedDomains *[]string  `json:"allowed_domains" validate:"omitempty,dive,fqdn"`
	Description    *string    `json:"description" validate:"omitempty,max=500"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}

// EventResponse represents the response body for an event
type EventResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	AllowedDomains []string  `json:"allowed_domains"`
	Description    *string   `json:"description,omitempty"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// EventListResponse represents the response body for a list of events
type EventListResponse struct {
	Events []EventResponse `json:"events"`
	Total  int             `json:"total"`
}
