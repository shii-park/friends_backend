package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/sqlc"
)

type GetUserService struct {
	queries *sqlc.Queries
}

func NewGetUserService(queries *sqlc.Queries) *GetUserService {
	return &GetUserService{queries: queries}
}

// GetUser はユーザーIDに対応するユーザー情報をDBから取得します
func (s *GetUserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	row, err := s.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           row.UserID.String(),
		Name:         row.UserName,
		Icon:         row.IconUrl.String,
		ProfileMessage: row.ProfileMessage.String,
		Birthmonth:   int(row.BirthMonth.Int16),
		Birthday:     int(row.BirthDay.Int16),
		RegisteredAt: row.RegisteredDate,
		StreakLogin:  int(row.StreakLoginDays),
		RankPoint:    int(row.RankPoint),
		Coin:         int(row.Coin),
		GachaStone:   int(row.GachaStone),
	}

	if row.LatestLoginDate.Valid {
		user.LatestLoginAt = row.LatestLoginDate.Time
	}

	return user, nil
}
