package refreshtoken

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/service"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrSessionExpired      = errors.New("session expired")
)

type Input struct {
	RefreshToken string
	IpAddress    string
	UserAgent    string
}

type Output struct {
	AccessToken  string
	RefreshToken string
}

type UseCase struct {
	jwtService  *service.JWTService
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewUseCase(
	jwtService *service.JWTService,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
) *UseCase {
	return &UseCase{
		jwtService:  jwtService,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
	// 1. Find session by refresh token
	session, err := uc.sessionRepo.FindByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidRefreshToken
	}

	// 2. Check if expired
	if session.IsExpired() {
		if err := uc.sessionRepo.Delete(ctx, session.ID); err != nil {
			return nil, err
		}
		return nil, ErrSessionExpired
	}

	// 3. Find user
	user, err := uc.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidRefreshToken
	}

	// 4. Generate new access token
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	// 5. Rotate refresh token (security)
	newRefreshToken, err := uc.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 6. Update session with new refresh token
	newSessionID, err := lib.GenerateID(lib.CUID2)
	if err != nil {
		return nil, err
	}

	// Delete old session and create new one (rotation)
	if err := uc.sessionRepo.Delete(ctx, session.ID); err != nil {
		return nil, err
	}

	newSession := entity.NewSession(
		newSessionID,
		user.ID,
		newRefreshToken,
		input.IpAddress,
		input.UserAgent,
		time.Now().Add(uc.jwtService.GetRefreshTokenTTL()),
	)

	if err := uc.sessionRepo.Save(ctx, newSession); err != nil {
		return nil, err
	}

	return &Output{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
