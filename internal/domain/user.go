package domain

import (
	"strings"
	"time"

	"friends_backend/internal/errs"
	"github.com/google/uuid"
)

const (
	ResetStreakLoginDay  = 1
	guestName = "ゲスト"
)


type User struct {
	ID             string    `json:"userID"` // 更新しない
	UserName       string    `json:"userName"`
	Icon           string    `json:"icon,omitempty"`// アイコン、stringにしてるが実際どうなるかはわからない
	ProfileMessage string    `json:"profileMsg,omitempty"`

	Birthday      time.Time `json:"birthday"` // 更新しない
	RegisteredAt  time.Time `json:"registeredDate"`
	LatestLoginAt time.Time `json:"latestLoginDate"`
	StreakLogin   int       `json:"streakLogin"`
}

// ユーザーIDを生成する関数
func NewUserID() string {
	return uuid.NewString()
}

// ユーザーを生成する関数
func NewUser(userName string, icon string, profileMessage string, birthday time.Time)  (*User, error) {
	name := strings.TrimSpace(userName)

	if name == "" {
		return nil, errs.ErrUserNameRequired
	}

	now := time.Now()

	return &User{
		ID:				NewUserID(),
		UserName:		name,
		Icon:			icon,
		ProfileMessage:	profileMessage,

		Birthday:		birthday,
		RegisteredAt:	now,
		LatestLoginAt:	now,
		StreakLogin:	1,
	}, nil
}

// ゲストユーザーを生成する関数
// ここにおくべきかはちょっと微妙ではある（一応置いておく）
func NewGuestUser() *User {
	user, _ := NewUser(guestName, "", "", time.Time{})
	return user
}

// ユーザー名の更新処理
func (u *User) UpdateUserName(userName string) error {
	name := strings.TrimSpace(userName)

	if name == "" {
		return errs.ErrUserNameRequired
	}

	if u.UserName == name {
		return nil
	}

	u.UserName = name

	return nil
}

// アイコンの更新処理
func (u *User) UpdateIcon(icon string){
	u.Icon = icon
}

// メッセージの更新処理
func (u *User) UpdateProfileMessage(profileMessage string){
	u.ProfileMessage = profileMessage
}

// ユーザーの更新処理
func (u *User) UpdateUser(userName string, icon string, profileMessage string) error {
	err := u.UpdateUserName(userName)

	if err != nil {
		return err
	}

	u.UpdateIcon(icon)
	u.UpdateProfileMessage(profileMessage)

	return nil
}

// ログイン処理
func (u *User) OnLogin() {
	prev := u.LatestLoginAt

	now := time.Now()
	u.LatestLoginAt = now

	if prev.IsZero() {
		u.ResetStreakLogin()
		return
	}

	prevDate := dateOnly(prev)
	nowDate := dateOnly(now)

	if prevDate.Equal(nowDate) {
		return
	}

	if prevDate.AddDate(0, 0, 1).Equal(nowDate) {
		u.AddStreakLogin()
	} else {
		u.ResetStreakLogin()
	}
}

// 日付だけ取得するメソッド
func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// 連続ログイン日数を増やす
func (u *User) AddStreakLogin(){
	u.StreakLogin++
}

// 連続ログイン日数のリセットする処理
func (u *User) ResetStreakLogin(){
	u.StreakLogin = ResetStreakLoginDay
}