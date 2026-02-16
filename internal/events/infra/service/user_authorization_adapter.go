package service

import (
	"context"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/service"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/repository"
)

type UserAuthorizationAdapter struct {
	userRepo repository.UserRepository
}

func NewUserAuthorizationAdapter(userRepo repository.UserRepository) *UserAuthorizationAdapter {
	return &UserAuthorizationAdapter{
		userRepo: userRepo,
	}
}

func (a *UserAuthorizationAdapter) GetUserByID(ctx context.Context, userID string) (*service.UserInfo, error) {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return &service.UserInfo{
		ID:        user.ID,
		IsAdmin:   user.IsAdmin(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

func (a *UserAuthorizationAdapter) IsUserAdmin(ctx context.Context, userID string) (bool, error) {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return user.IsAdmin(), nil
}

func (a *UserAuthorizationAdapter) GetUserEmail(ctx context.Context, userID string) (string, error) {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}
	return user.Email, nil
}

func (a *UserAuthorizationAdapter) GetUserInfoBatch(ctx context.Context, userIDs []string) ([]*service.UserInfo, error) {
	users, err := a.userRepo.FindByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*service.UserInfo, len(users))
	for i, user := range users {
		result[i] = &service.UserInfo{
			ID:        user.ID,
			IsAdmin:   user.IsAdmin(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
	}
	return result, nil
}
