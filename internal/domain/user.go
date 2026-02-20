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

func NewUserID() UserID {
	return uuid.NewString()
}


