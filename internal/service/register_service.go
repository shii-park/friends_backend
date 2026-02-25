package service

import (
	"context"
	"database/sql"

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

	err = s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		UserID:   user.ID,
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
		LatestLoginDate: sql.NullTime{
			Time:  user.LatestLoginAt,
			Valid: true,
		},
		StreakLoginDays: int32(user.StreakLogin),
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
