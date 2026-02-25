package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/errs"
	"github.com/shii-park/friends/internal/sqlc"
)

type storageRepository struct {
	q *sqlc.Queries
}

func NewStorageRepository(q *sqlc.Queries) domain.StorageRepository {
	return &storageRepository{q: q}
}

func (r *storageRepository) AddCard(
	ctx context.Context,
	userID uuid.UUID,
	instanceID uuid.UUID,
	cardID int,
) error {
	return r.q.AddUserCard(ctx, sqlc.AddUserCardParams{
		InstanceID: instanceID,
		UserID:     userID,
		CardID:     int32(cardID),
	})
}

func (r *storageRepository) RemoveCard(
	ctx context.Context,
	userID uuid.UUID,
	instanceID uuid.UUID,
) error {
	rows, err := r.q.RemoveUserCard(ctx, sqlc.RemoveUserCardParams{
		UserID:     userID,
		InstanceID: instanceID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrCardNotFound
	}
	return nil
}

func (r *storageRepository) ListCards(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.CardInstance, error) {
	rows, err := r.q.ListUserCards(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]domain.CardInstance, 0, len(rows))
	for _, row := range rows {
		res = append(res, domain.CardInstance{
			InstanceID: row.InstanceID,
			CardID:     int(row.CardID),
		})
	}

	return res, nil
}
