package mapper

import (
	createactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_activities"
	createevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_event"
	geteventwithactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_with_activities"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/http/dto"
)

func CreateEventRequestToInput(req *dto.CreateEventRequest) *createevent.Input {
	return &createevent.Input{
		Name:           req.Name,
		AllowedDomains: &req.AllowedDomains,
		Description:    req.Description,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	}
}

func EventToResponse(event *entity.Event) dto.EventResponse {
	return dto.EventResponse{
		ID:             event.ID,
		Name:           event.Name,
		AllowedDomains: event.AllowedDomains,
		Description:    event.Description,
		StartDate:      event.StartDate,
		EndDate:        event.EndDate,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func ActivityToResponse(activity *entity.Activity) dto.ActivityResponse {
	return dto.ActivityResponse{
		ID:          activity.ID,
		Name:        activity.Name,
		EventID:     activity.EventID,
		Description: activity.Description,
		StartDate:   activity.StartDate,
		EndDate:     activity.EndDate,
		CreatedAt:   activity.CreatedAt,
		UpdatedAt:   activity.UpdatedAt,
	}
}

func ActivitiesToResponse(activities []*entity.Activity) []dto.ActivityResponse {
	result := make([]dto.ActivityResponse, len(activities))
	for i, a := range activities {
		result[i] = ActivityToResponse(a)
	}
	return result
}

func GetEventWithActivitiesOutputToResponse(output *geteventwithactivities.Output) *dto.EventWithActivitiesResponse {
	return &dto.EventWithActivitiesResponse{
		Event:      EventToResponse(output.Event),
		Activities: ActivitiesToResponse(output.Activities),
	}
}

func CreateActivitiesRequestToInput(req *dto.CreateActivitiesRequest, eventID, userID string) *createactivities.Input {
	activities := make([]createactivities.ActivityInput, len(req.Activities))
	for i, a := range req.Activities {
		activities[i] = createactivities.ActivityInput{
			Name:        a.Name,
			Description: a.Description,
			StartDate:   a.StartDate,
			EndDate:     a.EndDate,
		}
	}
	return &createactivities.Input{
		UserID:     userID,
		EventID:    eventID,
		Activities: activities,
	}
}

func CreateActivitiesOutputToResponse(output *createactivities.Output) *dto.ActivitiesResponse {
	return &dto.ActivitiesResponse{
		Activities: ActivitiesToResponse(output.Activities),
	}
}
