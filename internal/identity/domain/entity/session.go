package entity

import "time"

type Session struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	IpAddress    string    `db:"ip_address"`
	UserAgent    string    `db:"user_agent"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
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
