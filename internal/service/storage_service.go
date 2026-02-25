package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
)

type StorageService interface {
	AddCard(ctx context.Context, userID uuid.UUID, cardID int) (domain.CardInstance, error)
	RemoveCard(ctx context.Context, userID uuid.UUID, instanceID uuid.UUID) error
	ListCards(ctx context.Context, userID uuid.UUID) ([]domain.CardInstance, error)

	ListCardDetails(ctx context.Context, userID uuid.UUID) ([]domain.CardInstanceDetail, error)
}

type storageService struct {
	repo domain.StorageRepository
}

func NewStorageService(repo domain.StorageRepository) StorageService {
	return &storageService{repo: repo}
}

func (s *storageService) AddCard(ctx context.Context, userID uuid.UUID, cardID int) (domain.CardInstance, error) {
	if userID == uuid.Nil {
		return domain.CardInstance{}, errs.ErrInvalidUserID
	}

	st := domain.NewEmptyStorage(userID)
	ci, err := st.Add(cardID)
	if err != nil {
		return domain.CardInstance{}, err
	}

	// DBに保存
	if err := s.repo.AddCard(ctx, userID, ci.InstanceID, ci.CardID); err != nil {
		return domain.CardInstance{}, err
	}

	return ci, nil
}

func (s *storageService) RemoveCard(ctx context.Context, userID uuid.UUID, instanceID uuid.UUID) error {
	if userID == uuid.Nil {
		return errs.ErrInvalidUserID
	}
	if instanceID == uuid.Nil {
		return errs.ErrInvalidInstanceID
	}
	return s.repo.RemoveCard(ctx, userID, instanceID)
}

func (s *storageService) ListCards(ctx context.Context, userID uuid.UUID) ([]domain.CardInstance, error) {
	if userID == uuid.Nil {
		return nil, errs.ErrInvalidUserID
	}
	return s.repo.ListCards(ctx, userID)
}

func (s *storageService) ListCardDetails(ctx context.Context, userID uuid.UUID) ([]domain.CardInstanceDetail, error) {
	if userID == uuid.Nil {
		return nil, errs.ErrInvalidUserID
	}
	return s.repo.ListCardDetails(ctx, userID)
}
