package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/domain"
	"github.com/shii-park/friends/internal/sqlc"
)

type RegisterService struct {
	queries *sqlc.Queries
}

func NewRegisterService(queries *sqlc.Queries) *RegisterService {
	return &RegisterService{queries: queries}
}

// RegisterUser はユーザーを生成してDBに保存します
func (s *RegisterService) RegisterUser(ctx context.Context, name string, icon string, profileMessage string, birthmonth int, birthday int) (*domain.User, error) {
	user, err := domain.NewUser(name, icon, profileMessage, birthmonth, birthday)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, err
	}

	err = s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		UserID:   userID,
		UserName: user.Name,
		IconUrl: sql.NullString{
			String: user.Icon,
			Valid:  user.Icon != "",
		},
		ProfileMessage: sql.NullString{
			String: user.ProfileMessage,
			Valid:  user.ProfileMessage != "",
		},
		BirthMonth: sql.NullInt16{
			Int16: int16(user.Birthmonth),
			Valid: true,
		},
		BirthDay: sql.NullInt16{
			Int16: int16(user.Birthday),
			Valid: true,
		},
		RegisteredDate: user.RegisteredAt,
		LatestLoginDate: sql.NullTime{
			Time:  user.LatestLoginAt,
			Valid: true,
		},
		StreakLoginDays: int32(user.StreakLogin),
		RankPoint:       int32(user.RankPoint),
		Coin:            int32(user.Coin),
		GachaStone:      int32(user.GachaStone),
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *RegisterService) DeleteUser(c context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return s.queries.DeleteUser(c, id)
}
