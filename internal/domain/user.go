package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shii-park/friends/internal/errs"
)

const (
	ResetStreakLoginDay = 1
	LoginCutoffHour     = 4

	guestName = "ゲスト"
)

type User struct {
	ID             string `json:"userID"` // 更新しない
	Name           string `json:"userName"`
	Icon           string `json:"icon,omitempty"` // アイコン、stringにしてるが実際どうなるかはわからない
	ProfileMessage string `json:"profileMsg,omitempty"`

	Birthday      time.Time `json:"birthday"` // 更新しない
	RegisteredAt  time.Time `json:"registeredDate"`
	LatestLoginAt time.Time `json:"latestLoginDate"`
	StreakLogin   int       `json:"streakLogin"`
}

// NewUserID ユーザーIDを生成する関数
func NewUserID() string {
	return uuid.NewString()
}

// NewUser ユーザーを生成する関数
func NewUser(name string, icon string, profileMessage string, birthday time.Time) (*User, error) {
	formattedName := strings.TrimSpace(name)

	if formattedName == "" {
		return nil, errs.UserNameRequired
	}

	now := time.Now()

	return &User{
		ID:             NewUserID(),
		Name:           formattedName,
		Icon:           icon,
		ProfileMessage: profileMessage,

		Birthday:      birthday,
		RegisteredAt:  now,
		LatestLoginAt: now,
		StreakLogin:   1,
	}, nil
}

// NewGuestUser ゲストユーザーを生成する関数
// ここにおくべきかはちょっと微妙ではある（一応置いておく）
func NewGuestUser() *User {
	user, _ := NewUser(guestName, "", "", time.Time{})
	return user
}

// UpdateUserName ユーザー名の更新処理
func (u *User) UpdateUserName(name string) error {
	formattedName := strings.TrimSpace(name)

	if name == "" {
		return errs.UserNameRequired
	}

	if u.Name == formattedName {
		return nil
	}

	u.Name = name

	return nil
}

// UpdateIcon アイコンの更新処理
func (u *User) UpdateIcon(icon string) {
	u.Icon = icon
}

// UpdateProfileMessage メッセージの更新処理
func (u *User) UpdateProfileMessage(profileMessage string) {
	u.ProfileMessage = profileMessage
}

// UpdateUser ユーザーの更新処理
func (u *User) UpdateUser(userName string, icon string, profileMessage string) error {
	err := u.UpdateUserName(userName)
	if err != nil {
		return err
	}

	u.UpdateIcon(icon)
	u.UpdateProfileMessage(profileMessage)

	return nil
}

// OnLogin ログイン処理
func (u *User) OnLogin() {
	prev := u.LatestLoginAt

	now := time.Now()
	u.LatestLoginAt = now

	if prev.IsZero() {
		u.ResetStreakLogin()
		return
	}

	prevDay := loginDay(prev)
	nowDay := loginDay(now)

	if prevDay.Equal(nowDay) {
		return
	}

	if prevDay.AddDate(0, 0, 1).Equal(nowDay) {
		u.AddStreakLogin()
	} else {
		u.ResetStreakLogin()
	}
}

// 日付だけ取り出して4時間前にする関数
func loginDay(t time.Time) time.Time {
	shifted := t.Add(-LoginCutoffHour * time.Hour)
	y, m, d := shifted.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, shifted.Location())
}

// AddStreakLogin 連続ログイン日数を増やす関数
func (u *User) AddStreakLogin() {
	u.StreakLogin++
}

// ResetStreakLogin 連続ログイン日数のリセットする関数
func (u *User) ResetStreakLogin() {
	u.StreakLogin = ResetStreakLoginDay
}
