package identity

import (
	"context"
	"errors"
	"time"

	identity "github.com/gabrielmatsan/checkin-gate/internal/identity/domain"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra"
	repository "github.com/gabrielmatsan/checkin-gate/internal/identity/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrSessionExpired      = errors.New("session expired")
)

type RefreshTokenInput struct {
	RefreshToken string
	IpAddress    string
	UserAgent    string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenUseCase struct {
	jwtService  *infra.JWTService
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewRefreshTokenUseCase(
	jwtService *infra.JWTService,
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		jwtService:  jwtService,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	// 1. Busca session pelo refresh token
	session, err := uc.sessionRepo.FindByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrInvalidRefreshToken
	}

	// 2. Verifica se expirou
	if session.IsExpired() {
		uc.sessionRepo.Delete(ctx, session.ID)
		return nil, ErrSessionExpired
	}

	// 3. Busca usuário
	user, err := uc.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidRefreshToken
	}

	// 4. Gera novo access token
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	// 5. Rotaciona refresh token (segurança)
	newRefreshToken, err := uc.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 6. Atualiza session com novo refresh token
	newSessionID, err := lib.GenerateID(lib.CUID2)
	if err != nil {
		return nil, err
	}

	// Deleta session antiga e cria nova (rotação)
	uc.sessionRepo.Delete(ctx, session.ID)

	newSession := identity.NewSession(
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

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
