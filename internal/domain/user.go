package domain

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID             string    `json:"userID"`
	UserName       string    `json:"userName"`
	Icon           string    `json:"icon,omitempty"`// アイコン、stringにしてるが実際どうなるかはわからない
	ProfileMessage string    `json:"profileMsg,omitempty"`

	Birthday      time.Time `json:"birthday"`
	RegisteredAt  time.Time `json:"registerdDate"`
	LatestLoginAt time.Time `json:"latestLoginDate"`
	StreakLogin   int       `json:"streakLogin"`
}

// ユーザーIDを生成するメソッド
func NewUserID() {
	return uuid.NewString()
}

// ユーザーを生成するメソッド
func NewUser(userName string, icon string, profileMessage string, birthday time.Time)  (*User, error) {
	name := strings.TrimSpace(userName)

	if name == "" {
		return nil, errors.New("userName is required")
	}

	now := time.Now()

	return &User{
		ID:				NewUserID(),
		UserName:		userName,
		Icon:			icon,
		ProfileMessage:	profileMessage,

		Birthday:		birthday,
		RegisteredAt:	now,
		LatestLoginAt:	now,
		StreakLogin:	1,
	}, nil
}

func NewGuestUser() (*User, error){
	user, err := NewUser("ゲスト", "", "", time.Time{})
	return user, err
}