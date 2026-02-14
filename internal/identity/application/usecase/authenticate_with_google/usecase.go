package authenticatewithgoogle

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/service"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type Input struct {
	Code      string
	IpAddress string
	UserAgent string
}

type Output struct {
	AccessToken  string
	RefreshToken string
	User         *entity.User
}

type UseCase struct {
	googleProvider *lib.GoogleOAuthProvider
	jwtService     *service.JWTService
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
}

func NewUseCase(
	googleProvider *lib.GoogleOAuthProvider,
	jwtService *service.JWTService,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
) *UseCase {
	return &UseCase{
		googleProvider: googleProvider,
		jwtService:     jwtService,
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
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
		user = entity.NewUser(entity.NewUserParams{
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

	session := entity.NewSession(
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

	return &Output{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
