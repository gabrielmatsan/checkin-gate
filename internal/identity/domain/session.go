package identity

import "time"

type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	IpAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

func NewSession(id, userID, refreshToken, ipAddress, userAgent string, expiresAt time.Time) *Session {
	return &Session{
		ID:           id,
		UserID:       userID,
		RefreshToken: refreshToken,
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
	}
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
