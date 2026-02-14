package listusercheckins

import (
	"context"
	"fmt"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
)

type Input struct {
	UserID string
}

type Output struct {
	CheckIns []*entity.CheckIn
}

type UseCase struct {
	checkInRepo repository.CheckInRepository
}

func NewUseCase(checkInRepo repository.CheckInRepository) *UseCase {
	return &UseCase{
		checkInRepo: checkInRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
	checkIns, err := uc.checkInRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user check-ins: %w", err)
	}

	return &Output{CheckIns: checkIns}, nil
}
