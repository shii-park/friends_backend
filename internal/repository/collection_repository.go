package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/sqlc"
)

type CollectionRepository interface {
	ListCollection(ctx context.Context, userID uuid.UUID) ([]sqlc.ListCollectionRow, error)
}

type collectionRepository struct {
	q *sqlc.Queries
}

func NewCollectionRepository(q *sqlc.Queries) CollectionRepository {
	return &collectionRepository{q: q}
}

func (r *collectionRepository) ListCollection(ctx context.Context, userID uuid.UUID) ([]sqlc.ListCollectionRow, error) {
	return r.q.ListCollection(ctx, userID)
}
