package identity

import (
	"fmt"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type GetGoogleAuthURLOutput struct {
	URL   string
	State string
}

type GetGoogleAuthURLUseCase struct {
	googleProvider *lib.GoogleOAuthProvider
}

func NewGetGoogleAuthURLUseCase(googleProvider *lib.GoogleOAuthProvider) *GetGoogleAuthURLUseCase {
	return &GetGoogleAuthURLUseCase{
		googleProvider: googleProvider,
	}
}

func (uc *GetGoogleAuthURLUseCase) Execute() (*GetGoogleAuthURLOutput, error) {
	state, err := lib.GenerateID(lib.CUID2)

	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	url := uc.googleProvider.GetAuthURL(state)

	return &GetGoogleAuthURLOutput{
		URL:   url,
		State: state,
	}, nil
}
