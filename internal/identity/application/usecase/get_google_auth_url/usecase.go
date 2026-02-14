package getgoogleauthurl

import (
	"fmt"

	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type Output struct {
	URL   string
	State string
}

type UseCase struct {
	googleProvider *lib.GoogleOAuthProvider
}

func NewUseCase(googleProvider *lib.GoogleOAuthProvider) *UseCase {
	return &UseCase{
		googleProvider: googleProvider,
	}
}

func (uc *UseCase) Execute() (*Output, error) {
	state, err := lib.GenerateID(lib.CUID2)
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	url := uc.googleProvider.GetAuthURL(state)

	return &Output{
		URL:   url,
		State: state,
	}, nil
}
