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
	ActivityID   string    `json:"activity_id"`
	ActivityName string    `json:"activity_name"`
	ActivityDate time.Time `json:"activity_date"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

type EventInfo struct {
	EventID   string `json:"event_id"`
	EventName string `json:"event_name"`
}

// certificateJob é a struct privada - só pode ser criada via NewCertificateJob
type certificateJob struct {
	JobID        string       `json:"job_id"`
	EventInfo    EventInfo    `json:"event_info"`
	UserInfo     UserInfo     `json:"user_info"`
	ActivityInfo ActivityInfo `json:"activity_info"`
	CheckedAt    time.Time    `json:"checked_at"`
	EnqueuedAt   time.Time    `json:"enqueued_at"`
}

// CertificateJob é o tipo público (alias) para uso externo
type CertificateJob = certificateJob

type NewCertificateJobParams struct {
	EventInfo    EventInfo
	UserInfo     UserInfo
	ActivityInfo ActivityInfo
	CheckedAt    time.Time
}

// NewCertificateJob é a única forma de criar um CertificateJob
func NewCertificateJob(params NewCertificateJobParams) (*CertificateJob, error) {
	id, err := lib.GenerateID(lib.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate job ID: %w", err)
	}

	return &CertificateJob{
		JobID:        id,
		EventInfo:    params.EventInfo,
		UserInfo:     params.UserInfo,
		ActivityInfo: params.ActivityInfo,
		CheckedAt:    params.CheckedAt,
		EnqueuedAt:   time.Now(),
	}, nil
}

// Getters para acesso aos campos (opcional, mas útil para encapsulamento)
func (j *CertificateJob) GetJobID() string        { return j.JobID }
func (j *CertificateJob) GetEventInfo() EventInfo { return j.EventInfo }
func (j *CertificateJob) GetUserInfo() UserInfo   { return j.UserInfo }
func (j *CertificateJob) GetActivityInfo() ActivityInfo { return j.ActivityInfo }
func (j *CertificateJob) GetCheckedAt() time.Time { return j.CheckedAt }
func (j *CertificateJob) GetEnqueuedAt() time.Time { return j.EnqueuedAt }
