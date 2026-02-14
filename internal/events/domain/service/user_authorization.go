package service

import "context"

type UserInfo struct {
	ID      string
	IsAdmin bool
}

type UserAuthorizationService interface {
	GetUserByID(ctx context.Context, userID string) (*UserInfo, error)
	IsUserAdmin(ctx context.Context, userID string) (bool, error)
}
