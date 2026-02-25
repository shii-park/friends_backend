package domain

import (
	"context"

	"github.com/google/uuid"
)

type StorageRepository interface {
	AddCard(
		ctx context.Context,
		userID uuid.UUID,
		instanceID uuid.UUID,
		cardID int,
	) error

	RemoveCard(
		ctx context.Context,
		userID uuid.UUID,
		instanceID uuid.UUID,
	) error

	ListCards(
		ctx context.Context,
		userID uuid.UUID,
	) ([]CardInstance, error)

	ListCardDetails(
		ctx context.Context,
		userID uuid.UUID,
	) ([]CardInstanceDetail, error)
}
