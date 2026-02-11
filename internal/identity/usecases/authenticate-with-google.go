package identity

import (
	"context"
	"fmt"
	"time"

	identity "github.com/gabrielmatsan/checkin-gate/internal/identity/domain"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra"
	repository "github.com/gabrielmatsan/checkin-gate/internal/identity/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type AuthenticateWithGoogleInput struct {
	Code      string `json:"code"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

type AuthenticateWithGoogleOutput struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	User         *identity.User `json:"user"`
}

type AuthenticateWithGoogleUseCase struct {
	googleProvider *lib.GoogleOAuthProvider
	jwtService     *infra.JWTService
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
}

func NewAuthenticateWithGoogleUseCase(googleProvider *lib.GoogleOAuthProvider, jwtService *infra.JWTService, userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *AuthenticateWithGoogleUseCase {
	return &AuthenticateWithGoogleUseCase{
		googleProvider: googleProvider,
		jwtService:     jwtService,
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
	}
}

func (uc *AuthenticateWithGoogleUseCase) Execute(ctx context.Context, input *AuthenticateWithGoogleInput) (*AuthenticateWithGoogleOutput, error) {

	token, err := uc.googleProvider.Enchange(ctx, input.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	getUserInfo, err := uc.googleProvider.GetUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := uc.userRepo.FindByEmail(ctx, getUserInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if user == nil {
		id, err := lib.GenerateID(lib.CUID2)
		if err != nil {
			return nil, fmt.Errorf("failed to generate user ID: %w", err)
		}
		user = identity.NewUser(identity.NewUserParams{
			ID:        id,
			FirstName: getUserInfo.FirstName,
			LastName:  getUserInfo.LastName,
			Email:     getUserInfo.Email,
		})

		user, err = uc.userRepo.Save(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to save user: %w", err)
		}
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	sessionID, err := lib.GenerateID(lib.CUID2)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	session := identity.NewSession(
		sessionID,
		user.ID,
		refreshToken,
		input.IpAddress,
		input.UserAgent,
		time.Now().Add(uc.jwtService.GetRefreshTokenTTL()),
	)

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return &AuthenticateWithGoogleOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
