package queue

import (
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type UserInfo struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

type ActivityInfo struct {
	ActivityID        string    `json:"activity_id"`
	ActivityName      string    `json:"activity_name"`
	ActivityDate      time.Time `json:"activity_date"`
	ActivityStartDate time.Time `json:"activity_start_date"`
	ActivityEndDate   time.Time `json:"activity_end_date"`
}

type EventInfo struct {
	EventID   string `json:"event_id"`
	EventName string `json:"event_name"`
}

type CertificateJob struct {
	JobID        string       `json:"job_id"`
	EventInfo    EventInfo    `json:"event_info"`
	UserInfo     UserInfo     `json:"user_info"`
	ActivityInfo ActivityInfo `json:"activity_info"`
}

type NewCertificateJobParams struct {
	EventInfo    EventInfo
	UserInfo     UserInfo
	ActivityInfo ActivityInfo
}

func NewCertificateJob(params NewCertificateJobParams) *CertificateJob {

	id, err := lib.GenerateID(lib.UUID)
	if err != nil {
		fmt.Errorf("failed to generate job ID: %w", err)
		return nil
	}

	return &CertificateJob{
		JobID:        id,
		EventInfo:    params.EventInfo,
		UserInfo:     params.UserInfo,
		ActivityInfo: params.ActivityInfo,
	}
}
