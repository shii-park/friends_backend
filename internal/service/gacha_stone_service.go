package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/repository"
)

type GachaStoneService struct {
	userRepo repository.UserRepository
}

func NewGachaStoneService(userRepo repository.UserRepository) *GachaStoneService {
	return &GachaStoneService{userRepo: userRepo}
}

func (s *GachaStoneService) Add(ctx context.Context, userID uuid.UUID, amount int) (int, error) {
	if amount <= 0 {
		return 0, errs.ErrInvalidGachaStoneDelta
	}
	return s.userRepo.AddGachaStone(ctx, userID, amount)
}
