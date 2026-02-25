package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/sqlc"
)

type UserRepository interface {
	AddGachaStone(ctx context.Context, userID uuid.UUID, amount int) (int, error)
}

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) AddGachaStone(ctx context.Context, userID uuid.UUID, amount int) (int, error) {
	v, err := r.q.AddUserGachaStone(ctx, sqlc.AddUserGachaStoneParams{
		UserID: userID,
		Amount: int32(amount),
	})
	if err != nil {
		return 0, err
	}
	return int(v), nil
}
