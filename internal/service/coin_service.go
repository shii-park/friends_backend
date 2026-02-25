package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/repository"
)

type CoinService struct {
	userRepo repository.UserRepository
}

func NewCoinService(userRepo repository.UserRepository) *CoinService {
	return &CoinService{userRepo: userRepo}
}

func (s *CoinService) Add(ctx context.Context, userID uuid.UUID, amount int) (int, error) {
	if amount <= 0 {
		return 0, errs.ErrInvalidCoinDelta
	}
	return s.userRepo.AddCoin(ctx, userID, amount)
}
