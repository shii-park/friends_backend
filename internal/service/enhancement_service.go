package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/repository"
)

type EnhancementService interface {
	EnhanceCard(ctx context.Context, userID uuid.UUID, instanceID uuid.UUID, times int) (newLevel int, cost int, err error)
}

type enhancementService struct {
	repo repository.EnhancementRepository
}

func NewEnhancementService(repo repository.EnhancementRepository) EnhancementService {
	return &enhancementService{repo: repo}
}

func (s *enhancementService) EnhanceCard(ctx context.Context, userID uuid.UUID, instanceID uuid.UUID, times int) (int, int, error) {
	return s.repo.EnhanceCardByCoin(ctx, userID, instanceID, times)
}
