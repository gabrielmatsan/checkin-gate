package dto

import "time"

type CreateEventRequest struct {
	Name           string    `json:"name" validate:"required,min=3,max=255"`
	AllowedDomains []string  `json:"allowed_domains" validate:"dive,fqdn"`
	Description    *string   `json:"description" validate:"omitempty,max=500"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateEventRequest struct {
	Name           *string    `json:"name" validate:"omitempty,min=3,max=255"`
	AllowedDomains *[]string  `json:"allowed_domains" validate:"omitempty,dive,fqdn"`
	Description    *string    `json:"description" validate:"omitempty,max=500"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}

type GetEventWithActivitiesRequest struct {
	EventID string `json:"event_id" validate:"required"`
}

type CreateActivityRequest struct {
	Name        string    `json:"name" validate:"required,min=3,max=255"`
	Description *string   `json:"description" validate:"omitempty,max=500"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type CreateActivitiesRequest struct {
	Activities []CreateActivityRequest `json:"activities" validate:"required,min=1,max=10,dive"`
}
