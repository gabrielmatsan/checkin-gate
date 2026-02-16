package service

import "context"

type UserInfo struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	IsAdmin   bool
}

type UserAuthorizationService interface {
	GetUserByID(ctx context.Context, userID string) (*UserInfo, error)
	IsUserAdmin(ctx context.Context, userID string) (bool, error)
	GetUserEmail(ctx context.Context, userID string) (string, error)
	GetUserInfoBatch(ctx context.Context, userIDs []string) ([]*UserInfo, error)
}
